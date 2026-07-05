import 'dart:async';
import 'package:flutter/material.dart';
import '../services/api_service.dart';

class SatuanBarangScreen extends StatefulWidget {
  final bool active;

  const SatuanBarangScreen({super.key, this.active = false});

  @override
  State<SatuanBarangScreen> createState() => _SatuanBarangScreenState();
}

class _SatuanBarangScreenState extends State<SatuanBarangScreen> {
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
  void didUpdateWidget(SatuanBarangScreen old) {
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
        if (event == 'satuan_barang' && mounted) _fetch();
      },
      onError: (_) {},
      cancelOnError: false,
    );
  }

  Future<void> _fetch() async {
    _hasFetched = true;
    setState(() => _loading = true);
    try {
      final res = await ApiService.get('/satuan-barang');
      setState(() => _list = res['data'] as List? ?? []);
    } catch (_) {}
    if (mounted) setState(() => _loading = false);
  }

  Future<void> _showForm({Map? item}) async {
    final ctrlSatuan = TextEditingController(text: item?['satuan'] ?? '');
    final ctrlKet = TextEditingController(text: item?['keterangan'] ?? '');
    final isEdit = item != null;

    final ok = await showDialog<bool>(
      context: context,
      builder: (ctx) => AlertDialog(
        title: Text(isEdit ? 'Edit Satuan' : 'Tambah Satuan'),
        content: Column(
          mainAxisSize: MainAxisSize.min,
          children: [
            TextField(
              controller: ctrlSatuan,
              decoration: const InputDecoration(labelText: 'Satuan', border: OutlineInputBorder()),
            ),
            const SizedBox(height: 12),
            TextField(
              controller: ctrlKet,
              decoration: const InputDecoration(labelText: 'Keterangan', border: OutlineInputBorder()),
            ),
          ],
        ),
        actions: [
          TextButton(onPressed: () => Navigator.pop(ctx), child: const Text('Batal')),
          ElevatedButton(
            onPressed: () async {
              try {
                final body = {
                  'satuan': ctrlSatuan.text.trim(),
                  'keterangan': ctrlKet.text.trim(),
                };
                if (isEdit) {
                  await ApiService.put('/satuan-barang/${item['id']}', body);
                } else {
                  await ApiService.post('/satuan-barang', body);
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

    ctrlSatuan.dispose();
    ctrlKet.dispose();

    if (ok == true) _fetch();
  }

  Future<void> _delete(int id) async {
    final ok = await showDialog<bool>(
      context: context,
      builder: (ctx) => AlertDialog(
        title: const Text('Hapus Satuan'),
        content: const Text('Yakin hapus?'),
        actions: [
          TextButton(onPressed: () => Navigator.pop(ctx, false), child: const Text('Batal')),
          TextButton(onPressed: () => Navigator.pop(ctx, true), child: const Text('Hapus')),
        ],
      ),
    );
    if (ok != true) return;
    try {
      await ApiService.delete('/satuan-barang/$id');
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
            label: const Text('Tambah Satuan'),
          ),
        ),
        Expanded(
          child: _list.isEmpty
              ? const Center(child: Text('Belum ada satuan'))
              : ListView.builder(
                  itemCount: _list.length,
                  itemBuilder: (_, i) {
                    final item = _list[i];
                    return Card(
                      margin: const EdgeInsets.symmetric(horizontal: 12, vertical: 4),
                      child: ListTile(
                        title: Text(item['satuan'] ?? ''),
                        subtitle: item['keterangan'] != null ? Text(item['keterangan']) : null,
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
