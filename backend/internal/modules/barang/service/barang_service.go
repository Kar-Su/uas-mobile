package service

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/Kar-Su/uas-mobile.git/internal/database/entities"
	"github.com/Kar-Su/uas-mobile.git/internal/modules/barang/dto"
	"github.com/Kar-Su/uas-mobile.git/internal/modules/barang/repository"
	satuanRepo "github.com/Kar-Su/uas-mobile.git/internal/modules/satuan_barang/repository"
	tipeRepo "github.com/Kar-Su/uas-mobile.git/internal/modules/tipe_barang/repository"
	"github.com/Kar-Su/uas-mobile.git/internal/package/constants"
	"github.com/Kar-Su/uas-mobile.git/internal/package/utils"

	"gorm.io/gorm"
)

type BarangService interface {
	GetAll(ctx context.Context, page, limit int, filter *dto.FilterBarangQuery) ([]dto.BarangResponse, int64, error)
	GetByKode(ctx context.Context, kode string) (dto.BarangResponse, error)
	Create(ctx context.Context, req dto.BarangCreateRequest) (dto.BarangResponse, error)
	Update(ctx context.Context, kode string, req dto.BarangUpdateRequest) (dto.BarangResponse, error)
	Delete(ctx context.Context, kode string) error
}

type barangService struct {
	repo       repository.BarangRepository
	tipeRepo   tipeRepo.TipeBarangRepository
	satuanRepo satuanRepo.SatuanBarangRepository
	db         *gorm.DB
}

func NewBarangService(
	repo repository.BarangRepository,
	tipeRepo tipeRepo.TipeBarangRepository,
	satuanRepo satuanRepo.SatuanBarangRepository,
	db *gorm.DB,
) BarangService {
	return &barangService{repo: repo, tipeRepo: tipeRepo, satuanRepo: satuanRepo, db: db}
}

func toResponse(b entities.Barang) dto.BarangResponse {
	return dto.BarangResponse{
		Kode:      b.Kode,
		Name:      b.Name,
		Quantity:  b.Quantity,
		CreatedAt: b.CreatedAt,
		UpdatedAt: b.UpdatedAt,
		Tipe:      dto.TipeResponse{ID: b.Tipe.ID, Name: b.Tipe.Name},
		Satuan:    dto.SatuanResponse{ID: b.Satuan.ID, Satuan: b.Satuan.Satuan, Keterangan: b.Satuan.Keterangan},
	}
}

func (s *barangService) GetAll(ctx context.Context, page, limit int, filter *dto.FilterBarangQuery) ([]dto.BarangResponse, int64, error) {
	pageSize := utils.ResolvePageSize(limit)
	offset := utils.GetOffsetWithSize(page, pageSize)
	data, total, err := s.repo.FindAll(ctx, s.db, offset, pageSize, filter)
	if err != nil {
		log.Printf("Internal Error: %v", err)
		return nil, 0, constants.ErrInternalErr
	}
	result := make([]dto.BarangResponse, len(data))
	for i, d := range data {
		result[i] = toResponse(d)
	}
	return result, total, nil
}

func (s *barangService) GetByKode(ctx context.Context, kode string) (dto.BarangResponse, error) {
	data, err := s.repo.FindByKode(ctx, s.db, kode)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return dto.BarangResponse{}, dto.ErrBarangNotFound
		}
		log.Printf("Internal Error: %v", err)
		return dto.BarangResponse{}, constants.ErrInternalErr
	}
	return toResponse(data), nil
}

func (s *barangService) Create(ctx context.Context, req dto.BarangCreateRequest) (dto.BarangResponse, error) {
	if req.Quantity < 0 {
		return dto.BarangResponse{}, errors.New("quantity cannot be negative")
	}
	if _, err := s.tipeRepo.FindByID(ctx, s.db, req.TipeID); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return dto.BarangResponse{}, errors.New("tipe barang not found")
		}
		log.Printf("Internal Error: %v", err)
		return dto.BarangResponse{}, constants.ErrInternalErr
	}
	if _, err := s.satuanRepo.FindByID(ctx, s.db, req.SatuanID); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return dto.BarangResponse{}, errors.New("satuan barang not found")
		}
		log.Printf("Internal Error: %v", err)
		return dto.BarangResponse{}, constants.ErrInternalErr
	}
	now := time.Now()
	data, err := s.repo.Create(ctx, s.db, entities.Barang{
		Kode:      req.Kode,
		Name:      req.Name,
		TipeID:    req.TipeID,
		SatuanID:  req.SatuanID,
		Quantity:  req.Quantity,
		CreatedAt: now,
		UpdatedAt: now,
	})
	if err != nil {
		log.Printf("Internal Error: %v", err)
		return dto.BarangResponse{}, constants.ErrInternalErr
	}
	return toResponse(data), nil
}

func (s *barangService) Update(ctx context.Context, kode string, req dto.BarangUpdateRequest) (dto.BarangResponse, error) {
	if _, err := s.repo.FindByKode(ctx, s.db, kode); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return dto.BarangResponse{}, dto.ErrBarangNotFound
		}
		log.Printf("Internal Error: %v", err)
		return dto.BarangResponse{}, constants.ErrInternalErr
	}
	patch := entities.Barang{UpdatedAt: time.Now()}
	if req.Name != "" {
		patch.Name = req.Name
	}
	if req.TipeID != 0 {
		if _, err := s.tipeRepo.FindByID(ctx, s.db, req.TipeID); err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return dto.BarangResponse{}, errors.New("tipe barang not found")
			}
			log.Printf("Internal Error: %v", err)
			return dto.BarangResponse{}, constants.ErrInternalErr
		}
		patch.TipeID = req.TipeID
	}
	if req.SatuanID != 0 {
		if _, err := s.satuanRepo.FindByID(ctx, s.db, req.SatuanID); err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return dto.BarangResponse{}, errors.New("satuan barang not found")
			}
			log.Printf("Internal Error: %v", err)
			return dto.BarangResponse{}, constants.ErrInternalErr
		}
		patch.SatuanID = req.SatuanID
	}
	if req.Quantity != nil {
		if *req.Quantity < 0 {
			return dto.BarangResponse{}, errors.New("quantity cannot be negative")
		}
		patch.Quantity = *req.Quantity
	}
	data, err := s.repo.Update(ctx, s.db, kode, patch)
	if err != nil {
		log.Printf("Internal Error: %v", err)
		return dto.BarangResponse{}, constants.ErrInternalErr
	}
	return toResponse(data), nil
}

func (s *barangService) Delete(ctx context.Context, kode string) error {
	if _, err := s.repo.FindByKode(ctx, s.db, kode); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return dto.ErrBarangNotFound
		}
		log.Printf("Internal Error: %v", err)
		return constants.ErrInternalErr
	}
	if err := s.repo.Delete(ctx, s.db, kode); err != nil {
		log.Printf("Internal Error: %v", err)
		return constants.ErrInternalErr
	}
	return nil
}
