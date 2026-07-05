import 'dart:async';
import 'package:flutter/material.dart';
import '../services/api_service.dart';

class TipeBarangScreen extends StatefulWidget {
  final bool active;

  const TipeBarangScreen({super.key, this.active = false});

  @override
  State<TipeBarangScreen> createState() => _TipeBarangScreenState();
}

class _TipeBarangScreenState extends State<TipeBarangScreen> {
  List _list = [];
  bool _loading = true;
  bool _hasFetched = false;
  StreamSubscription<String>? _sseSub;

  @override
  void initState() {
    super.initState();
    if (widget.active) {
      _fetch();
      _startSSE();
    }
  }

  @override
  void didUpdateWidget(TipeBarangScreen old) {
    super.didUpdateWidget(old);
    if (widget.active && !_hasFetched) {
      _fetch();
      _startSSE();
    }
  }

  @override
  void dispose() {
    _sseSub?.cancel();
    super.dispose();
  }

  void _startSSE() {
    _sseSub = ApiService.sseStream('/sse').listen(
      (event) {
        if (event == 'tipe_barang' && mounted) _fetch();
      },
      onError: (_) {},
      cancelOnError: false,
    );
  }

  Future<void> _fetch() async {
    _hasFetched = true;
    setState(() => _loading = true);
    try {
      final res = await ApiService.get('/tipe-barang');
      setState(() => _list = res['data'] as List? ?? []);
    } catch (_) {}
    if (mounted) setState(() => _loading = false);
  }

  Future<void> _showForm({Map? item}) async {
    final ctrl = TextEditingController(text: item?['name'] ?? '');
    final isEdit = item != null;

    final ok = await showDialog<bool>(
      context: context,
      builder: (ctx) => AlertDialog(
        title: Text(isEdit ? 'Edit Tipe' : 'Tambah Tipe'),
        content: TextField(
          controller: ctrl,
          decoration: const InputDecoration(labelText: 'Nama Tipe', border: OutlineInputBorder()),
        ),
        actions: [
          TextButton(onPressed: () => Navigator.pop(ctx), child: const Text('Batal')),
          ElevatedButton(
            onPressed: () async {
              try {
                if (isEdit) {
                  await ApiService.put('/tipe-barang/${item['id']}', {'name': ctrl.text.trim()});
                } else {
                  await ApiService.post('/tipe-barang', {'name': ctrl.text.trim()});
                }
                if (ctx.mounted) Navigator.pop(ctx, true);
              } catch (e) {
                ScaffoldMessenger.of(ctx).showSnackBar(SnackBar(content: Text('Gagal: $e')));
              }
            },
            child: Text(isEdit ? 'Simpan' : 'Tambah'),
          ),
        ],
      ),
    );

    ctrl.dispose();

    if (ok == true) _fetch();
  }

  Future<void> _delete(int id) async {
    final ok = await showDialog<bool>(
      context: context,
      builder: (ctx) => AlertDialog(
        title: const Text('Hapus Tipe'),
        content: const Text('Yakin hapus?'),
        actions: [
          TextButton(onPressed: () => Navigator.pop(ctx, false), child: const Text('Batal')),
          TextButton(onPressed: () => Navigator.pop(ctx, true), child: const Text('Hapus')),
        ],
      ),
    );
    if (ok != true) return;
    try {
      await ApiService.delete('/tipe-barang/$id');
      _fetch();
    } catch (e) {
      if (mounted) ScaffoldMessenger.of(context).showSnackBar(SnackBar(content: Text('Gagal: $e')));
    }
  }

  @override
  Widget build(BuildContext context) {
    if (_loading) return const Center(child: CircularProgressIndicator());
    return Column(
      children: [
        Padding(
          padding: const EdgeInsets.all(8),
          child: ElevatedButton.icon(
            onPressed: () => _showForm(),
            icon: const Icon(Icons.add),
            label: const Text('Tambah Tipe'),
          ),
        ),
        Expanded(
          child: _list.isEmpty
              ? const Center(child: Text('Belum ada tipe'))
              : ListView.builder(
                  itemCount: _list.length,
                  itemBuilder: (_, i) {
                    final item = _list[i];
                    return Card(
                      margin: const EdgeInsets.symmetric(horizontal: 12, vertical: 4),
                      child: ListTile(
                        title: Text(item['name'] ?? ''),
                        trailing: Row(
                          mainAxisSize: MainAxisSize.min,
                          children: [
                            IconButton(icon: const Icon(Icons.edit), onPressed: () => _showForm(item: item)),
                            IconButton(icon: const Icon(Icons.delete), onPressed: () => _delete(item['id'])),
                          ],
                        ),
                      ),
                    );
                  },
                ),
        ),
      ],
    );
  }
}
