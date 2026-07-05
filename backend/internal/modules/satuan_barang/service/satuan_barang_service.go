package service

import (
	"context"
	"errors"
	"log"

	"github.com/Kar-Su/uas-mobile.git/internal/database/entities"
	"github.com/Kar-Su/uas-mobile.git/internal/modules/satuan_barang/dto"
	"github.com/Kar-Su/uas-mobile.git/internal/modules/satuan_barang/repository"
	"github.com/Kar-Su/uas-mobile.git/internal/package/constants"
	"github.com/Kar-Su/uas-mobile.git/internal/package/utils"

	"gorm.io/gorm"
)

type SatuanBarangService interface {
	GetAll(ctx context.Context) ([]dto.SatuanBarangResponse, error)
	GetByID(ctx context.Context, id int) (dto.SatuanBarangResponse, error)
	Create(ctx context.Context, req dto.SatuanBarangRequest) (dto.SatuanBarangResponse, error)
	Update(ctx context.Context, id int, req dto.SatuanBarangRequest) (dto.SatuanBarangResponse, error)
	Delete(ctx context.Context, id int) error
}

type satuanBarangService struct {
	repo repository.SatuanBarangRepository
	db   *gorm.DB
}

func NewSatuanBarangService(repo repository.SatuanBarangRepository, db *gorm.DB) SatuanBarangService {
	return &satuanBarangService{repo: repo, db: db}
}

func toResponse(s entities.SatuanBarang) dto.SatuanBarangResponse {
	return dto.SatuanBarangResponse{ID: s.ID, Satuan: s.Satuan, Keterangan: s.Keterangan}
}

func (s *satuanBarangService) GetAll(ctx context.Context) ([]dto.SatuanBarangResponse, error) {
	data, err := s.repo.FindAll(ctx, s.db)
	if err != nil {
		log.Printf("Internal Error: %v", err)
		return nil, constants.ErrInternalErr
	}
	result := make([]dto.SatuanBarangResponse, len(data))
	for i, d := range data {
		result[i] = toResponse(d)
	}
	return result, nil
}

func (s *satuanBarangService) GetByID(ctx context.Context, id int) (dto.SatuanBarangResponse, error) {
	data, err := s.repo.FindByID(ctx, s.db, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return dto.SatuanBarangResponse{}, dto.ErrSatuanBarangNotFound
		}
		log.Printf("Internal Error: %v", err)
		return dto.SatuanBarangResponse{}, constants.ErrInternalErr
	}
	return toResponse(data), nil
}

func (s *satuanBarangService) Create(ctx context.Context, req dto.SatuanBarangRequest) (dto.SatuanBarangResponse, error) {
	satuan := req.Satuan
	keterangan := req.Keterangan
	data, err := s.repo.Create(ctx, s.db, entities.SatuanBarang{
		Satuan:     &satuan,
		Keterangan: &keterangan,
	})
	if err != nil {
		log.Printf("Internal Error: %v", err)
		return dto.SatuanBarangResponse{}, constants.ErrInternalErr
	}
	return toResponse(data), nil
}

func (s *satuanBarangService) Update(ctx context.Context, id int, req dto.SatuanBarangRequest) (dto.SatuanBarangResponse, error) {
	_, err := s.repo.FindByID(ctx, s.db, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return dto.SatuanBarangResponse{}, dto.ErrSatuanBarangNotFound
		}
		log.Printf("Internal Error: %v", err)
		return dto.SatuanBarangResponse{}, constants.ErrInternalErr
	}
	satuan := req.Satuan
	keterangan := req.Keterangan
	data, err := s.repo.Update(ctx, s.db, id, entities.SatuanBarang{
		Satuan:     &satuan,
		Keterangan: &keterangan,
	})
	if err != nil {
		log.Printf("Internal Error: %v", err)
		return dto.SatuanBarangResponse{}, constants.ErrInternalErr
	}
	return toResponse(data), nil
}

func (s *satuanBarangService) Delete(ctx context.Context, id int) error {
	_, err := s.repo.FindByID(ctx, s.db, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return dto.ErrSatuanBarangNotFound
		}
		log.Printf("Internal Error: %v", err)
		return constants.ErrInternalErr
	}
	if err := s.repo.Delete(ctx, s.db, id); err != nil {
		if utils.IsForeignKeyViolation(err) {
			return dto.ErrSatuanBarangInUse
		}
		log.Printf("Internal Error: %v", err)
		return constants.ErrInternalErr
	}
	return nil
}
