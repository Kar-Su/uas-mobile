import 'dart:async';
import 'dart:convert';
import 'package:http/http.dart' as http;
import 'package:shared_preferences/shared_preferences.dart';
import '../config.dart';

class ApiService {
  static const String _baseUrl = AppConfig.baseUrl;
  static const Duration _timeout = Duration(seconds: 15);

  static Future<http.Response> _request(Future<http.Response> Function() fn) {
    return fn().timeout(_timeout);
  }

  static const String _tokenKey = 'access_token';
  static const String _refreshKey = 'refresh_token';
  static const String _roleKey = 'role_name';
  static const String _nameKey = 'user_name';

  static Future<String?> getAccessToken() async {
    final prefs = await SharedPreferences.getInstance();
    return prefs.getString(_tokenKey);
  }

  static Future<void> saveTokens({
    required String accessToken,
    required String refreshToken,
    required String roleName,
    required String userName,
  }) async {
    final prefs = await SharedPreferences.getInstance();
    await prefs.setString(_tokenKey, accessToken);
    await prefs.setString(_refreshKey, refreshToken);
    await prefs.setString(_roleKey, roleName);
    await prefs.setString(_nameKey, userName);
  }

  static Future<void> clearTokens() async {
    final prefs = await SharedPreferences.getInstance();
    await prefs.remove(_tokenKey);
    await prefs.remove(_refreshKey);
    await prefs.remove(_roleKey);
    await prefs.remove(_nameKey);
  }

  static Future<String?> getRole() async {
    final prefs = await SharedPreferences.getInstance();
    return prefs.getString(_roleKey);
  }

  static Future<String?> getUserName() async {
    final prefs = await SharedPreferences.getInstance();
    return prefs.getString(_nameKey);
  }

  static Future<Map<String, String>> _headers() async {
    final token = await getAccessToken();
    return {
      'Content-Type': 'application/json',
      if (token != null) 'Authorization': 'Bearer $token',
    };
  }

  static Future<Map<String, dynamic>> _refreshToken() async {
    final prefs = await SharedPreferences.getInstance();
    final refreshToken = prefs.getString(_refreshKey);
    if (refreshToken == null) throw Exception('No refresh token');

    final res = await _request(
      () => http.post(
        Uri.parse('$_baseUrl/auth/refresh-token'),
        headers: {'Content-Type': 'application/json'},
        body: jsonEncode({'refresh_token': refreshToken}),
      ),
    );

    if (res.statusCode != 200 && res.statusCode != 201) {
      throw Exception('Refresh failed');
    }

    final body = jsonDecode(res.body);
    final data = body['data'];
    _extractCookiesToData(res, data);
    return data;
  }

  static Future<Map<String, dynamic>> get(
    String path, {
    bool allowRetry = true,
  }) async {
    final res = await _request(
      () async =>
          http.get(Uri.parse('$_baseUrl$path'), headers: await _headers()),
    );

    if (res.statusCode == 401 && allowRetry) {
      return _retryWithRefresh(() => get(path, allowRetry: false));
    }

    return _handleResponse(res);
  }

  static Future<Map<String, dynamic>> post(
    String path,
    Map<String, dynamic> body, {
    bool allowRetry = true,
  }) async {
    final res = await _request(
      () async => http.post(
        Uri.parse('$_baseUrl$path'),
        headers: await _headers(),
        body: jsonEncode(body),
      ),
    );

    if (res.statusCode == 401 && allowRetry) {
      return _retryWithRefresh(() => post(path, body, allowRetry: false));
    }

    return _handleResponse(res);
  }

  static Future<Map<String, dynamic>> put(
    String path,
    Map<String, dynamic> body, {
    bool allowRetry = true,
  }) async {
    final res = await _request(
      () async => http.put(
        Uri.parse('$_baseUrl$path'),
        headers: await _headers(),
        body: jsonEncode(body),
      ),
    );

    if (res.statusCode == 401 && allowRetry) {
      return _retryWithRefresh(() => put(path, body, allowRetry: false));
    }

    return _handleResponse(res);
  }

  static Future<Map<String, dynamic>> delete(
    String path, {
    bool allowRetry = true,
  }) async {
    final res = await _request(
      () async =>
          http.delete(Uri.parse('$_baseUrl$path'), headers: await _headers()),
    );

    if (res.statusCode == 401 && allowRetry) {
      return _retryWithRefresh(() => delete(path, allowRetry: false));
    }

    return _handleResponse(res);
  }

  static Future<Map<String, dynamic>> login(
    String email,
    String password,
  ) async {
    final res = await _request(
      () => http.post(
        Uri.parse('$_baseUrl/auth/login'),
        headers: {'Content-Type': 'application/json'},
        body: jsonEncode({'email': email, 'password': password}),
      ),
    );

    final responseMap = _handleResponse(res);
    if (responseMap['data'] != null) {
      _extractCookiesToData(res, responseMap['data']);
    }
    return responseMap;
  }

  static void _extractCookiesToData(http.Response res, Map<String, dynamic> data) {
    final cookieHeader = res.headers['set-cookie'];
    if (cookieHeader != null) {
      final parts = cookieHeader.split(';');
      for (final p in parts) {
        final strings = p.split(',');
        for (final s in strings) {
          final trimmed = s.trim();
          if (trimmed.startsWith('access_token=')) {
            data['access_token'] = trimmed.substring('access_token='.length);
          } else if (trimmed.startsWith('refresh_token=')) {
            data['refresh_token'] = trimmed.substring('refresh_token='.length);
          }
        }
      }
    }
  }

  static Future<void> logout() async {
    try {
      await post('/auth/logout', {});
    } catch (_) {}
    await clearTokens();
  }

  static Future<Map<String, dynamic>> _retryWithRefresh(
    Future<Map<String, dynamic>> Function() originalRequest,
  ) async {
    final data = await _refreshToken();
    await saveTokens(
      accessToken: data['access_token'],
      refreshToken: data['refresh_token'],
      roleName: data['role_name'],
      userName: data['user_name'],
    );
    return originalRequest();
  }

  static Map<String, dynamic> _handleResponse(http.Response res) {
    final body = jsonDecode(res.body) as Map<String, dynamic>;

    if (res.statusCode >= 200 && res.statusCode < 300) {
      return body;
    }

    final error = body['error'] ?? 'Unknown error';
    throw Exception(error.toString());
  }

  static Stream<String> sseStream(String path) async* {
    final token = await getAccessToken();
    final client = http.Client();
    try {
      final request = http.Request('GET', Uri.parse('$_baseUrl$path'));
      request.headers['Accept'] = 'text/event-stream';
      request.headers['Cache-Control'] = 'no-cache';
      if (token != null) request.headers['Authorization'] = 'Bearer $token';

      final response = await client.send(request);
      final stream = response.stream.transform(utf8.decoder).transform(const LineSplitter());

      await for (final line in stream) {
        if (line.startsWith('event:')) {
          yield line.substring(6).trim();
        }
      }
    } finally {
      client.close();
    }
  }
}
