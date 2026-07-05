import 'package:flutter/material.dart';
import '../services/api_service.dart';
import 'barang_list_screen.dart';
import 'tipe_barang_screen.dart';
import 'satuan_barang_screen.dart';
import 'login_screen.dart';

class AdminDashboard extends StatefulWidget {
  const AdminDashboard({super.key});

  @override
  State<AdminDashboard> createState() => _AdminDashboardState();
}

class _AdminDashboardState extends State<AdminDashboard> {
  int _pageIndex = 0;
  String _userName = '';

  @override
  void initState() {
    super.initState();
    _loadUser();
  }

  Future<void> _loadUser() async {
    final name = await ApiService.getUserName();
    if (mounted) setState(() => _userName = name ?? 'Admin');
  }

  Future<void> _logout() async {
    await ApiService.logout();
    if (!mounted) return;
    Navigator.pushReplacement(
      context,
      MaterialPageRoute(builder: (_) => const LoginScreen()),
    );
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        title: ListTile(
          title: Text('Dashboard Admin Gudang'),
          subtitle: Text('Selamat datang, $_userName'),
        ),

        actions: [
          IconButton(icon: const Icon(Icons.logout), onPressed: _logout),
        ],
      ),
      body: IndexedStack(
        index: _pageIndex,
        children: [
          BarangListScreen(active: _pageIndex == 0),
          TipeBarangScreen(active: _pageIndex == 1),
          SatuanBarangScreen(active: _pageIndex == 2),
        ],
      ),
      bottomNavigationBar: NavigationBar(
        selectedIndex: _pageIndex,
        onDestinationSelected: (i) => setState(() => _pageIndex = i),
        destinations: const [
          NavigationDestination(icon: Icon(Icons.inventory), label: 'Barang'),
          NavigationDestination(icon: Icon(Icons.category), label: 'Tipe'),
          NavigationDestination(icon: Icon(Icons.straighten), label: 'Satuan'),
        ],
      ),
    );
  }
}
