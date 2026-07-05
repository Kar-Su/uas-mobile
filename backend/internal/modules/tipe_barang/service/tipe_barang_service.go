package service

import (
	"context"
	"errors"
	"log"

	"github.com/Kar-Su/uas-mobile.git/internal/database/entities"
	"github.com/Kar-Su/uas-mobile.git/internal/modules/tipe_barang/dto"
	"github.com/Kar-Su/uas-mobile.git/internal/modules/tipe_barang/repository"
	"github.com/Kar-Su/uas-mobile.git/internal/package/constants"
	"github.com/Kar-Su/uas-mobile.git/internal/package/utils"

	"gorm.io/gorm"
)

type TipeBarangService interface {
	GetAll(ctx context.Context) ([]dto.TipeBarangResponse, error)
	GetByID(ctx context.Context, id int) (dto.TipeBarangResponse, error)
	Create(ctx context.Context, req dto.TipeBarangRequest) (dto.TipeBarangResponse, error)
	Update(ctx context.Context, id int, req dto.TipeBarangRequest) (dto.TipeBarangResponse, error)
	Delete(ctx context.Context, id int) error
}

type tipeBarangService struct {
	repo repository.TipeBarangRepository
	db   *gorm.DB
}

func NewTipeBarangService(repo repository.TipeBarangRepository, db *gorm.DB) TipeBarangService {
	return &tipeBarangService{repo: repo, db: db}
}

func toResponse(t entities.TipeBarang) dto.TipeBarangResponse {
	return dto.TipeBarangResponse{ID: t.ID, Name: t.Name}
}

func (s *tipeBarangService) GetAll(ctx context.Context) ([]dto.TipeBarangResponse, error) {
	data, err := s.repo.FindAll(ctx, s.db)
	if err != nil {
		log.Printf("Internal Error: %v", err)
		return nil, constants.ErrInternalErr
	}
	result := make([]dto.TipeBarangResponse, len(data))
	for i, d := range data {
		result[i] = toResponse(d)
	}
	return result, nil
}

func (s *tipeBarangService) GetByID(ctx context.Context, id int) (dto.TipeBarangResponse, error) {
	data, err := s.repo.FindByID(ctx, s.db, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return dto.TipeBarangResponse{}, dto.ErrTipeBarangNotFound
		}
		log.Printf("Internal Error: %v", err)
		return dto.TipeBarangResponse{}, constants.ErrInternalErr
	}
	return toResponse(data), nil
}

func (s *tipeBarangService) Create(ctx context.Context, req dto.TipeBarangRequest) (dto.TipeBarangResponse, error) {
	data, err := s.repo.Create(ctx, s.db, entities.TipeBarang{Name: req.Name})
	if err != nil {
		log.Printf("Internal Error: %v", err)
		return dto.TipeBarangResponse{}, constants.ErrInternalErr
	}
	return toResponse(data), nil
}

func (s *tipeBarangService) Update(ctx context.Context, id int, req dto.TipeBarangRequest) (dto.TipeBarangResponse, error) {
	_, err := s.repo.FindByID(ctx, s.db, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return dto.TipeBarangResponse{}, dto.ErrTipeBarangNotFound
		}
		log.Printf("Internal Error: %v", err)
		return dto.TipeBarangResponse{}, constants.ErrInternalErr
	}
	data, err := s.repo.Update(ctx, s.db, id, entities.TipeBarang{Name: req.Name})
	if err != nil {
		log.Printf("Internal Error: %v", err)
		return dto.TipeBarangResponse{}, constants.ErrInternalErr
	}
	return toResponse(data), nil
}

func (s *tipeBarangService) Delete(ctx context.Context, id int) error {
	_, err := s.repo.FindByID(ctx, s.db, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return dto.ErrTipeBarangNotFound
		}
		log.Printf("Internal Error: %v", err)
		return constants.ErrInternalErr
	}
	if err := s.repo.Delete(ctx, s.db, id); err != nil {
		if utils.IsForeignKeyViolation(err) {
			return dto.ErrTipeBarangInUse
		}
		log.Printf("Internal Error: %v", err)
		return constants.ErrInternalErr
	}
	return nil
}
