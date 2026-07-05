package repository

import (
	"context"

	"github.com/Kar-Su/uas-mobile.git/internal/database/entities"

	"gorm.io/gorm"
)

type SatuanBarangRepository interface {
	FindAll(ctx context.Context, tx *gorm.DB) ([]entities.SatuanBarang, error)
	FindByID(ctx context.Context, tx *gorm.DB, id int) (entities.SatuanBarang, error)
	Create(ctx context.Context, tx *gorm.DB, data entities.SatuanBarang) (entities.SatuanBarang, error)
	Update(ctx context.Context, tx *gorm.DB, id int, data entities.SatuanBarang) (entities.SatuanBarang, error)
	Delete(ctx context.Context, tx *gorm.DB, id int) error
}

type satuanBarangRepository struct {
	db *gorm.DB
}

func NewSatuanBarangRepository(db *gorm.DB) SatuanBarangRepository {
	return &satuanBarangRepository{db: db}
}

func (r *satuanBarangRepository) FindAll(ctx context.Context, tx *gorm.DB) ([]entities.SatuanBarang, error) {
	if tx == nil {
		tx = r.db
	}
	var result []entities.SatuanBarang
	if err := tx.WithContext(ctx).Find(&result).Error; err != nil {
		return nil, err
	}
	return result, nil
}

func (r *satuanBarangRepository) FindByID(ctx context.Context, tx *gorm.DB, id int) (entities.SatuanBarang, error) {
	if tx == nil {
		tx = r.db
	}
	var result entities.SatuanBarang
	if err := tx.WithContext(ctx).First(&result, id).Error; err != nil {
		return entities.SatuanBarang{}, err
	}
	return result, nil
}

func (r *satuanBarangRepository) Create(ctx context.Context, tx *gorm.DB, data entities.SatuanBarang) (entities.SatuanBarang, error) {
	if tx == nil {
		tx = r.db
	}
	if err := tx.WithContext(ctx).Create(&data).Error; err != nil {
		return entities.SatuanBarang{}, err
	}
	return data, nil
}

func (r *satuanBarangRepository) Update(ctx context.Context, tx *gorm.DB, id int, data entities.SatuanBarang) (entities.SatuanBarang, error) {
	if tx == nil {
		tx = r.db
	}
	if err := tx.WithContext(ctx).Model(&entities.SatuanBarang{}).Where("id = ?", id).Updates(map[string]any{
		"satuan":     data.Satuan,
		"keterangan": data.Keterangan,
	}).Error; err != nil {
		return entities.SatuanBarang{}, err
	}
	return r.FindByID(ctx, tx, id)
}

func (r *satuanBarangRepository) Delete(ctx context.Context, tx *gorm.DB, id int) error {
	if tx == nil {
		tx = r.db
	}
	if err := tx.WithContext(ctx).Delete(&entities.SatuanBarang{}, id).Error; err != nil {
		return err
	}
	return nil
}
