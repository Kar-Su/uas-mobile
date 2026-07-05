package entities

type Role struct {
	ID   int    `gorm:"primaryKey;autoIncrement"`
	Name string `gorm:"type:text;uniqueIndex;not null"`
}

func (Role) TableName() string { return "roles" }
