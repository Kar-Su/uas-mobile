package entities

type SatuanBarang struct {
	ID         int     `gorm:"primaryKey;autoIncrement"`
	Satuan     *string `gorm:"column:satuan;type:text"`
	Keterangan *string `gorm:"column:keterangan;type:text"`
}

func (SatuanBarang) TableName() string { return "satuan_barang" }
