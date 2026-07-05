package repository

import (
	"context"

	"github.com/Kar-Su/uas-mobile.git/internal/database/entities"

	"gorm.io/gorm"
)

type TipeBarangRepository interface {
	FindAll(ctx context.Context, tx *gorm.DB) ([]entities.TipeBarang, error)
	FindByID(ctx context.Context, tx *gorm.DB, id int) (entities.TipeBarang, error)
	Create(ctx context.Context, tx *gorm.DB, data entities.TipeBarang) (entities.TipeBarang, error)
	Update(ctx context.Context, tx *gorm.DB, id int, data entities.TipeBarang) (entities.TipeBarang, error)
	Delete(ctx context.Context, tx *gorm.DB, id int) error
}

type tipeBarangRepository struct {
	db *gorm.DB
}

func NewTipeBarangRepository(db *gorm.DB) TipeBarangRepository {
	return &tipeBarangRepository{db: db}
}

func (r *tipeBarangRepository) FindAll(ctx context.Context, tx *gorm.DB) ([]entities.TipeBarang, error) {
	if tx == nil {
		tx = r.db
	}
	var result []entities.TipeBarang
	if err := tx.WithContext(ctx).Find(&result).Error; err != nil {
		return nil, err
	}
	return result, nil
}

func (r *tipeBarangRepository) FindByID(ctx context.Context, tx *gorm.DB, id int) (entities.TipeBarang, error) {
	if tx == nil {
		tx = r.db
	}
	var result entities.TipeBarang
	if err := tx.WithContext(ctx).First(&result, id).Error; err != nil {
		return entities.TipeBarang{}, err
	}
	return result, nil
}

func (r *tipeBarangRepository) Create(ctx context.Context, tx *gorm.DB, data entities.TipeBarang) (entities.TipeBarang, error) {
	if tx == nil {
		tx = r.db
	}
	if err := tx.WithContext(ctx).Create(&data).Error; err != nil {
		return entities.TipeBarang{}, err
	}
	return data, nil
}

func (r *tipeBarangRepository) Update(ctx context.Context, tx *gorm.DB, id int, data entities.TipeBarang) (entities.TipeBarang, error) {
	if tx == nil {
		tx = r.db
	}
	if err := tx.WithContext(ctx).Model(&entities.TipeBarang{}).Where("id = ?", id).Updates(&data).Error; err != nil {
		return entities.TipeBarang{}, err
	}
	return r.FindByID(ctx, tx, id)
}

func (r *tipeBarangRepository) Delete(ctx context.Context, tx *gorm.DB, id int) error {
	if tx == nil {
		tx = r.db
	}
	if err := tx.WithContext(ctx).Delete(&entities.TipeBarang{}, id).Error; err != nil {
		return err
	}
	return nil
}
