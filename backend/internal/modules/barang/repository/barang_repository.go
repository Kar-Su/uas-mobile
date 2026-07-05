package repository

import (
	"context"

	"github.com/Kar-Su/uas-mobile.git/internal/database/entities"
	"github.com/Kar-Su/uas-mobile.git/internal/modules/barang/dto"
	"github.com/Kar-Su/uas-mobile.git/internal/package/utils"

	"gorm.io/gorm"
)

type BarangRepository interface {
	FindAll(ctx context.Context, tx *gorm.DB, offset, limit int, filter *dto.FilterBarangQuery) ([]entities.Barang, int64, error)
	FindByKode(ctx context.Context, tx *gorm.DB, kode string) (entities.Barang, error)
	Create(ctx context.Context, tx *gorm.DB, data entities.Barang) (entities.Barang, error)
	Update(ctx context.Context, tx *gorm.DB, kode string, data entities.Barang) (entities.Barang, error)
	Delete(ctx context.Context, tx *gorm.DB, kode string) error
}

type barangRepository struct {
	db *gorm.DB
}

func NewBarangRepository(db *gorm.DB) BarangRepository {
	return &barangRepository{db: db}
}

func (r *barangRepository) FindAll(ctx context.Context, tx *gorm.DB, offset, limit int, filter *dto.FilterBarangQuery) ([]entities.Barang, int64, error) {
	if tx == nil {
		tx = r.db
	}

	var result []entities.Barang

	query := tx.WithContext(ctx).Model(&entities.Barang{}).Preload("Tipe").Preload("Satuan")
	filterEmpty := utils.IsNilStruct(filter)

	if !filterEmpty {
		if filter.Search != nil {
			filter.Search.BindingQuery(query, []string{"name", "kode"})
		}
		if filter.Tipe != nil {
			query = query.Where("tipe_id = ?", *filter.Tipe)
		}
		if filter.QtyOrder != nil {
			query = query.Order("quantity " + *filter.QtyOrder)
		}
	}

	var total int64
	if err := query.Session(&gorm.Session{}).Select("kode").Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := query.Offset(offset).Limit(limit).Find(&result).Error; err != nil {
		return nil, 0, err
	}

	return result, total, nil
}

func (r *barangRepository) FindByKode(ctx context.Context, tx *gorm.DB, kode string) (entities.Barang, error) {
	if tx == nil {
		tx = r.db
	}
	var result entities.Barang
	if err := tx.WithContext(ctx).Preload("Tipe").Preload("Satuan").First(&result, "kode = ?", kode).Error; err != nil {
		return entities.Barang{}, err
	}
	return result, nil
}

func (r *barangRepository) Create(ctx context.Context, tx *gorm.DB, data entities.Barang) (entities.Barang, error) {
	if tx == nil {
		tx = r.db
	}
	if err := tx.WithContext(ctx).Create(&data).Error; err != nil {
		return entities.Barang{}, err
	}
	return r.FindByKode(ctx, tx, data.Kode)
}

func (r *barangRepository) Update(ctx context.Context, tx *gorm.DB, kode string, data entities.Barang) (entities.Barang, error) {
	if tx == nil {
		tx = r.db
	}
	if err := tx.WithContext(ctx).Model(&entities.Barang{}).Where("kode = ?", kode).Updates(&data).Error; err != nil {
		return entities.Barang{}, err
	}
	return r.FindByKode(ctx, tx, kode)
}

func (r *barangRepository) Delete(ctx context.Context, tx *gorm.DB, kode string) error {
	if tx == nil {
		tx = r.db
	}
	if err := tx.WithContext(ctx).Where("kode = ?", kode).Delete(&entities.Barang{}).Error; err != nil {
		return err
	}
	return nil
}
