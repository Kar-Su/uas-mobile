class Barang {
  final String kode;
  final String name;
  final int quantity;
  final TipeInfo tipe;
  final SatuanInfo satuan;

  Barang({
    required this.kode,
    required this.name,
    required this.quantity,
    required this.tipe,
    required this.satuan,
  });

  factory Barang.fromJson(Map<String, dynamic> json) {
    return Barang(
      kode: json['kode'] ?? '',
      name: json['name'] ?? '',
      quantity: json['quantity'] ?? 0,
      tipe: TipeInfo.fromJson(json['tipe'] ?? {}),
      satuan: SatuanInfo.fromJson(json['satuan'] ?? {}),
    );
  }
}

class TipeInfo {
  final int id;
  final String name;

  TipeInfo({required this.id, required this.name});

  factory TipeInfo.fromJson(Map<String, dynamic> json) {
    return TipeInfo(
      id: json['id'] ?? 0,
      name: json['name'] ?? '',
    );
  }
}

class SatuanInfo {
  final int id;
  final String? satuan;

  SatuanInfo({required this.id, this.satuan});

  factory SatuanInfo.fromJson(Map<String, dynamic> json) {
    return SatuanInfo(
      id: json['id'] ?? 0,
      satuan: json['satuan'],
    );
  }
}

class TipeBarang {
  final int id;
  final String name;

  TipeBarang({required this.id, required this.name});

  factory TipeBarang.fromJson(Map<String, dynamic> json) {
    return TipeBarang(
      id: json['id'] ?? 0,
      name: json['name'] ?? '',
    );
  }
}

class SatuanBarang {
  final int id;
  final String? satuan;
  final String? keterangan;

  SatuanBarang({required this.id, this.satuan, this.keterangan});

  factory SatuanBarang.fromJson(Map<String, dynamic> json) {
    return SatuanBarang(
      id: json['id'] ?? 0,
      satuan: json['satuan'],
      keterangan: json['keterangan'],
    );
  }
}
