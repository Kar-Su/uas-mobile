import 'dart:async';
import 'package:flutter/material.dart';
import '../models/barang.dart';
import '../services/api_service.dart';

class BarangListScreen extends StatefulWidget {
  final bool readOnly;
  final VoidCallback? onChanged;
  final bool active;

  const BarangListScreen({
    super.key,
    this.readOnly = false,
    this.onChanged,
    this.active = true,
  });

  @override
  State<BarangListScreen> createState() => _BarangListScreenState();
}

class _BarangListScreenState extends State<BarangListScreen> {
  List<Barang> _list = [];
  bool _loading = true;
  String? _error;
  int _page = 1;
  int _totalPages = 1;

  String _search = '';
  int _tipeFilter = 0;
  String _qtyOrder = '';
  List<Map<String, dynamic>> _tipeOptions = [];

  final TextEditingController _searchCtrl = TextEditingController();
  StreamSubscription<String>? _sseSub;

  @override
  void initState() {
    super.initState();
    _fetchOptions();
    _fetch();
    _startSSE();
  }

  @override
  void dispose() {
    _sseSub?.cancel();
    _searchCtrl.dispose();
    super.dispose();
  }

  void _startSSE() {
    _sseSub = ApiService.sseStream('/sse').listen(
      (event) {
        if (event == 'barang' && mounted) _fetch();
      },
      onError: (_) {},
      cancelOnError: false,
    );
  }

  Future<void> _fetchOptions() async {
    try {
      final tipeRes = await ApiService.get('/tipe-barang');
      setState(() {
        _tipeOptions = (tipeRes['data'] as List? ?? [])
            .cast<Map<String, dynamic>>();
      });
    } catch (_) {}
  }

  Future<void> _fetch() async {
    if (!mounted) return;
    setState(() {
      _loading = true;
      _error = null;
    });

    List<Barang> newList = [];
    int newTotalPages = 1;
    String? newError;

    try {
      final List<String> queryParams = [];
      if (_tipeFilter > 0) {
        queryParams.add('tipe=${_tipeFilter.toString()}');
      }
      if (_search.isNotEmpty) {
        queryParams.add('search=$_search');
      }
      if (_qtyOrder.isNotEmpty) {
        queryParams.add('qty_order=$_qtyOrder');
      }
      final res = await ApiService.get(
        '/barang?page=$_page${queryParams.isNotEmpty ? '&' : ''}${queryParams.join('&')}',
      );
      final raw = res['data'] as List? ?? [];
      final meta = res['meta'] as Map<String, dynamic>?;
      newList = raw.map((e) => Barang.fromJson(e)).toList();
      newTotalPages = meta?['total_pages'] as int? ?? 1;
    } catch (e) {
      newError = e.toString();
    }

    if (mounted) {
      setState(() {
        _list = newList;
        _totalPages = newTotalPages;
        _error = newError;
        _loading = false;
      });
    }
  }

  Future<void> _resetAndFetch() async {
    _searchCtrl.clear();
    setState(() {
      _search = '';
      _tipeFilter = 0;
      _qtyOrder = '';
      _page = 1;
    });
    await _fetch();
  }

  void _onSearchChanged(String v) {
    setState(() => _search = v);
  }

  void _onSearchSubmitted(String v) {
    setState(() => _search = v);
    _page = 1;
    _fetch();
  }

  Future<void> _goToPage(int page) async {
    if (page < 1 || page > _totalPages) return;
    setState(() => _page = page);
    await _fetch();
  }

  void _onTipeChanged(int? v) {
    setState(() {
      _tipeFilter = v ?? 0;
      _page = 1;
    });
    _fetch();
  }

  void _onSortPressed() {
    setState(() {
      if (_qtyOrder == '') {
        _qtyOrder = 'asc';
      } else if (_qtyOrder == 'asc') {
        _qtyOrder = 'desc';
      } else {
        _qtyOrder = '';
      }
      _page = 1;
    });
    _fetch();
  }

  Future<void> _delete(String kode) async {
    final ok = await showDialog<bool>(
      context: context,
      builder: (ctx) => AlertDialog(
        title: const Text('Hapus Barang'),
        content: Text('Yakin hapus barang $kode?'),
        actions: [
          TextButton(
            onPressed: () => Navigator.pop(ctx, false),
            child: const Text('Batal'),
          ),
          TextButton(
            onPressed: () => Navigator.pop(ctx, true),
            child: const Text('Hapus'),
          ),
        ],
      ),
    );
    if (ok != true) return;

    try {
      await ApiService.delete('/barang/$kode');
      _fetch();
      widget.onChanged?.call();
    } catch (e) {
      if (mounted) {
        ScaffoldMessenger.of(
          context,
        ).showSnackBar(SnackBar(content: Text('Gagal hapus: $e')));
      }
    }
  }

  Future<void> _showForm({Barang? barang}) async {
    final ctrlKode = TextEditingController(text: barang?.kode ?? '');
    final ctrlName = TextEditingController(text: barang?.name ?? '');
    final ctrlQty = TextEditingController(
      text: barang?.quantity.toString() ?? '0',
    );
    int? tipeId = barang?.tipe.id;
    int? satuanId = barang?.satuan.id;

    List tipeList = [];
    List satuanList = [];
    try {
      final tipeRes = await ApiService.get('/tipe-barang');
      tipeList = tipeRes['data'] as List? ?? [];
      final satuanRes = await ApiService.get('/satuan-barang');
      satuanList = satuanRes['data'] as List? ?? [];
    } catch (_) {}

    final isEdit = barang != null;

    final result = await showDialog<bool>(
      context: context,
      builder: (ctx) => StatefulBuilder(
        builder: (ctx, setDialogState) {
          return AlertDialog(
            title: Text(isEdit ? 'Edit Barang' : 'Tambah Barang'),
            contentPadding: const EdgeInsets.fromLTRB(20, 12, 20, 4),
            content: SingleChildScrollView(
              child: Column(
                mainAxisSize: MainAxisSize.min,
                children: [
                  TextField(
                    controller: ctrlKode,
                    decoration: const InputDecoration(
                      labelText: 'Kode',
                      border: OutlineInputBorder(),
                      isDense: true,
                      contentPadding: EdgeInsets.symmetric(
                        vertical: 10,
                        horizontal: 12,
                      ),
                    ),
                    enabled: !isEdit,
                  ),
                  const SizedBox(height: 10),
                  TextField(
                    controller: ctrlName,
                    decoration: const InputDecoration(
                      labelText: 'Nama',
                      border: OutlineInputBorder(),
                      isDense: true,
                      contentPadding: EdgeInsets.symmetric(
                        vertical: 10,
                        horizontal: 12,
                      ),
                    ),
                  ),
                  const SizedBox(height: 10),
                  TextField(
                    controller: ctrlQty,
                    decoration: const InputDecoration(
                      labelText: 'Quantity',
                      border: OutlineInputBorder(),
                      isDense: true,
                      contentPadding: EdgeInsets.symmetric(
                        vertical: 10,
                        horizontal: 12,
                      ),
                    ),
                    keyboardType: TextInputType.number,
                  ),
                  const SizedBox(height: 10),
                  DropdownButtonFormField<int>(
                    value: tipeId,
                    isExpanded: true,
                    dropdownColor: Colors.white,
                    decoration: const InputDecoration(
                      labelText: 'Tipe',
                      border: OutlineInputBorder(),
                      isDense: true,
                      contentPadding: EdgeInsets.symmetric(
                        vertical: 10,
                        horizontal: 12,
                      ),
                    ),
                    style: TextStyle(color: Colors.black),
                    icon: Icon(Icons.arrow_drop_down, color: Colors.black54),
                    items: tipeList
                        .map<DropdownMenuItem<int>>(
                          (e) => DropdownMenuItem<int>(
                            value: e['id'] as int,
                            child: Text(
                              e['name'] ?? '',
                              style: TextStyle(color: Colors.black),
                            ),
                          ),
                        )
                        .toList(),
                    onChanged: (v) => setDialogState(() => tipeId = v),
                  ),
                  const SizedBox(height: 10),
                  DropdownButtonFormField<int>(
                    value: satuanId,
                    isExpanded: true,
                    dropdownColor: Colors.white,
                    decoration: const InputDecoration(
                      labelText: 'Satuan',
                      border: OutlineInputBorder(),
                      isDense: true,
                      contentPadding: EdgeInsets.symmetric(
                        vertical: 10,
                        horizontal: 12,
                      ),
                    ),
                    style: TextStyle(color: Colors.black),
                    icon: Icon(Icons.arrow_drop_down, color: Colors.black54),
                    items: satuanList
                        .map<DropdownMenuItem<int>>(
                          (e) => DropdownMenuItem<int>(
                            value: e['id'] as int,
                            child: Text(
                              e['satuan'] ?? '',
                              style: TextStyle(color: Colors.black),
                            ),
                          ),
                        )
                        .toList(),
                    onChanged: (v) => setDialogState(() => satuanId = v),
                  ),
                ],
              ),
            ),
            actions: [
              TextButton(
                onPressed: () => Navigator.pop(ctx, false),
                child: const Text('Batal'),
              ),
              ElevatedButton(
                onPressed: () async {
                  try {
                    final body = {
                      'kode': ctrlKode.text.trim(),
                      'name': ctrlName.text.trim(),
                      'quantity': int.tryParse(ctrlQty.text) ?? 0,
                      'tipe_barang_id': tipeId,
                      'satuan_barang_id': satuanId,
                    };
                    if (isEdit) {
                      await ApiService.put('/barang/${barang.kode}', body);
                    } else {
                      await ApiService.post('/barang', body);
                    }
                    if (ctx.mounted) Navigator.pop(ctx, true);
                  } catch (e) {
                    ScaffoldMessenger.of(
                      ctx,
                    ).showSnackBar(SnackBar(content: Text('Gagal: $e')));
                  }
                },
                child: Text(isEdit ? 'Simpan' : 'Tambah'),
              ),
            ],
          );
        },
      ),
    );

    ctrlKode.dispose();
    ctrlName.dispose();
    ctrlQty.dispose();

    if (result == true) {
      _fetch();
      widget.onChanged?.call();
    }
  }

  @override
  Widget build(BuildContext context) {
    if (_loading) return const Center(child: CircularProgressIndicator());
    if (_error != null) return Center(child: Text('Error: $_error'));

    return Column(
      children: [
        Padding(
          padding: const EdgeInsets.fromLTRB(12, 8, 12, 4),
          child: Row(
            children: [
              Expanded(
                flex: 3,
                child: TextField(
                  controller: _searchCtrl,
                  onChanged: _onSearchChanged,
                  onSubmitted: _onSearchSubmitted,
                  textInputAction: TextInputAction.search,
                  decoration: InputDecoration(
                    hintText: 'Search...',
                    prefixIcon: const Icon(Icons.search, size: 18),
                    isDense: true,
                    contentPadding: const EdgeInsets.symmetric(
                      vertical: 8,
                      horizontal: 8,
                    ),
                    border: OutlineInputBorder(
                      borderRadius: BorderRadius.circular(8),
                    ),
                  ),
                  style: const TextStyle(fontSize: 14),
                ),
              ),
              const SizedBox(width: 6),
              Expanded(
                flex: 2,
                child: DropdownButtonFormField<int>(
                  value: _tipeFilter == 0 ? null : _tipeFilter,
                  isExpanded: true,
                  dropdownColor: Colors.white,
                  decoration: InputDecoration(
                    hintText: 'All Tipe',
                    hintStyle: const TextStyle(
                      color: Colors.black54,
                      fontSize: 13,
                    ),
                    isDense: true,
                    contentPadding: const EdgeInsets.symmetric(
                      vertical: 8,
                      horizontal: 8,
                    ),
                    border: OutlineInputBorder(
                      borderRadius: BorderRadius.circular(8),
                    ),
                  ),
                  style: const TextStyle(fontSize: 13, color: Colors.black),
                  icon: const Icon(
                    Icons.arrow_drop_down,
                    color: Colors.black54,
                  ),
                  items: [
                    const DropdownMenuItem<int>(
                      value: 0,
                      child: Text(
                        'All Tipe',
                        style: TextStyle(fontSize: 13, color: Colors.black),
                      ),
                    ),
                    ..._tipeOptions.map((e) {
                      return DropdownMenuItem<int>(
                        value: e['id'] as int,
                        child: Text(
                          e['name'] ?? '',
                          style: const TextStyle(
                            fontSize: 13,
                            color: Colors.black,
                          ),
                        ),
                      );
                    }),
                  ],
                  onChanged: _onTipeChanged,
                ),
              ),
              const SizedBox(width: 2),
              IconButton(
                constraints: const BoxConstraints(minWidth: 36, minHeight: 36),
                icon: Icon(
                  _qtyOrder == 'asc'
                      ? Icons.arrow_upward
                      : _qtyOrder == 'desc'
                      ? Icons.arrow_downward
                      : Icons.swap_vert,
                  size: 18,
                ),
                onPressed: _onSortPressed,
                tooltip: 'Sort by quantity',
                padding: EdgeInsets.zero,
              ),
            ],
          ),
        ),
        if (!widget.readOnly)
          Padding(
            padding: const EdgeInsets.all(8),
            child: ElevatedButton.icon(
              onPressed: () => _showForm(),
              icon: const Icon(Icons.add),
              label: const Text('Tambah Barang'),
            ),
          ),
        Expanded(
          child: RefreshIndicator(
            onRefresh: _resetAndFetch,
            child: _list.isEmpty
                ? const Center(child: Text('Belum ada barang'))
                : ListView.builder(
                    itemCount: _list.length,
                    itemBuilder: (_, i) {
                      final b = _list[i];
                      return Card(
                        margin: const EdgeInsets.symmetric(
                          horizontal: 12,
                          vertical: 4,
                        ),
                        child: ListTile(
                          title: Text(b.name),
                          subtitle: Text(
                            '${b.kode} | ${b.quantity} ${b.satuan.satuan ?? ""} | ${b.tipe.name}',
                          ),
                          trailing: widget.readOnly
                              ? null
                              : Row(
                                  mainAxisSize: MainAxisSize.min,
                                  children: [
                                    IconButton(
                                      icon: const Icon(Icons.edit),
                                      onPressed: () => _showForm(barang: b),
                                    ),
                                    IconButton(
                                      icon: const Icon(Icons.delete),
                                      onPressed: () => _delete(b.kode),
                                    ),
                                  ],
                                ),
                        ),
                      );
                    },
                  ),
          ),
        ),
        if (_totalPages > 1)
          Padding(
            padding: const EdgeInsets.only(bottom: 8),
            child: Row(
              mainAxisAlignment: MainAxisAlignment.center,
              children: [
                IconButton(
                  icon: const Icon(Icons.chevron_left),
                  onPressed: _page > 1 ? () => _goToPage(_page - 1) : null,
                ),
                Text('$_page / $_totalPages'),
                IconButton(
                  icon: const Icon(Icons.chevron_right),
                  onPressed: _page < _totalPages
                      ? () => _goToPage(_page + 1)
                      : null,
                ),
              ],
            ),
          ),
      ],
    );
  }
}
