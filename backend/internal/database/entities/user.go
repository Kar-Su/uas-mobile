package entities

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	ID       uuid.UUID `gorm:"type:uuid;primaryKey"`
	Email    string    `gorm:"column:email;type:citext;not null"`
	Name     string    `gorm:"column:name;type:text;not null"`
	Password string    `gorm:"column:password;type:text;not null"`
	RoleID   int       `gorm:"column:role_id;not null"`
	Role     Role      `gorm:"foreignKey:RoleID"`
}

func (u *User) BeforeCreate(tx *gorm.DB) error {
	if u.ID == uuid.Nil {
		u.ID = uuid.New()
	}
	return nil
}

func (User) TableName() string { return "users" }
