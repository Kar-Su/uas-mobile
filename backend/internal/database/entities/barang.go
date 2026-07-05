package entities

import "time"

type Barang struct {
	Kode      string       `gorm:"column:kode;type:text;primaryKey"`
	Name      string       `gorm:"column:name;type:text;not null"`
	TipeID    int          `gorm:"column:tipe_id;not null"`
	SatuanID  int          `gorm:"column:satuan_id;not null"`
	Quantity  int          `gorm:"column:quantity;default:0"`
	CreatedAt time.Time    `gorm:"column:created_at;not null"`
	UpdatedAt time.Time    `gorm:"column:updated_at;not null"`
	Tipe      TipeBarang   `gorm:"foreignKey:TipeID"`
	Satuan    SatuanBarang `gorm:"foreignKey:SatuanID"`
}

func (Barang) TableName() string { return "barang" }
