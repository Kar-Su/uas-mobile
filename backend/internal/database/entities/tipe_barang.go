package entities

type TipeBarang struct {
	ID   int    `gorm:"primaryKey;autoIncrement"`
	Name string `gorm:"column:name;type:text;not null"`
}

func (TipeBarang) TableName() string { return "tipe_barang" }
