package repository

import (
	"context"

	"github.com/Kar-Su/uas-mobile.git/internal/database/entities"
	"gorm.io/gorm"
)

type RoleRepository interface {
	FindByName(ctx context.Context, tx *gorm.DB, name string) (entities.Role, error)
	FindAll(ctx context.Context, tx *gorm.DB) ([]entities.Role, error)
}

type roleRepository struct {
	db *gorm.DB
}

func NewRoleRepository(db *gorm.DB) RoleRepository {
	return &roleRepository{db: db}
}

func (r *roleRepository) FindByName(ctx context.Context, tx *gorm.DB, name string) (entities.Role, error) {
	if tx == nil {
		tx = r.db
	}
	var role entities.Role
	if err := tx.WithContext(ctx).Where("name = ?", name).First(&role).Error; err != nil {
		return entities.Role{}, err
	}
	return role, nil
}

func (r *roleRepository) FindAll(ctx context.Context, tx *gorm.DB) ([]entities.Role, error) {
	if tx == nil {
		tx = r.db
	}
	var roles []entities.Role
	if err := tx.WithContext(ctx).Find(&roles).Error; err != nil {
		return nil, err
	}
	return roles, nil
}
