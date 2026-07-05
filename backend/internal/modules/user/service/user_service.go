package service

import (
	"context"
	"errors"
	"log"

	"github.com/Kar-Su/uas-mobile.git/internal/database/entities"
	"github.com/Kar-Su/uas-mobile.git/internal/modules/user/dto"
	"github.com/Kar-Su/uas-mobile.git/internal/modules/user/repository"
	"github.com/Kar-Su/uas-mobile.git/internal/package/constants"
	"github.com/Kar-Su/uas-mobile.git/internal/package/helpers"
	"github.com/Kar-Su/uas-mobile.git/internal/package/utils"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserService interface {
	GetAllUsers(ctx context.Context, page int, searchQuery *dto.SearchFilterQuery) ([]dto.UserResponse, int64, error)
	GetUserByID(ctx context.Context, id uuid.UUID) (dto.UserResponse, error)
	CreateUser(ctx context.Context, req dto.UserCreateRequest) error
	UpdateUser(ctx context.Context, id uuid.UUID, req dto.UserUpdateRequest) (dto.UserResponse, error)
	DeleteUser(ctx context.Context, id uuid.UUID) error
}

type userService struct {
	userRepo repository.UserRepository
	roleRepo repository.RoleRepository
	db       *gorm.DB
}

func NewUserService(userRepo repository.UserRepository, roleRepo repository.RoleRepository, db *gorm.DB) UserService {
	return &userService{
		userRepo: userRepo,
		roleRepo: roleRepo,
		db:       db,
	}
}

func toResponse(user entities.User) dto.UserResponse {
	return dto.UserResponse{
		ID:       user.ID,
		Email:    user.Email,
		Name:     user.Name,
		RoleName: user.Role.Name,
	}
}

func (s *userService) GetAllUsers(ctx context.Context, page int, searchFilterQuery *dto.SearchFilterQuery) ([]dto.UserResponse, int64, error) {
	if searchFilterQuery != nil {
		if searchFilterQuery.FilterRole != nil {
			*searchFilterQuery.FilterRole = helpers.NormalizeString(*searchFilterQuery.FilterRole)
		}
	}
	offset := utils.GetOffset(page)
	users, total, err := s.userRepo.FindAll(ctx, s.db, offset, utils.DefaultPageSize, searchFilterQuery)
	if err != nil {
		log.Printf("Internal Error: %v", err)
		return nil, 0, constants.ErrInternalErr
	}

	responses := make([]dto.UserResponse, len(users))
	for i, u := range users {
		u.Role.Name = helpers.NormalizeToResponseString(u.Role.Name)
		responses[i] = toResponse(u)
	}
	return responses, total, nil
}

func (s *userService) GetUserByID(ctx context.Context, id uuid.UUID) (dto.UserResponse, error) {
	user, err := s.userRepo.FindByID(ctx, s.db, id)
	user.Role.Name = helpers.NormalizeToResponseString(user.Role.Name)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return dto.UserResponse{}, dto.ErrUserNotFound
		}
		log.Printf("Internal Error: %v", err)
		return dto.UserResponse{}, constants.ErrInternalErr
	}
	return toResponse(user), nil
}

func (s *userService) CreateUser(ctx context.Context, req dto.UserCreateRequest) error {
	role, err := s.roleRepo.FindByName(ctx, s.db, req.RoleName)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("role not found")
		}
		log.Printf("Internal Error: %v", err)
		return constants.ErrInternalErr
	}

	if len(req.Password) < 8 {
		return errors.New("password minimal 8 karakter")
	}

	hashedPassword, err := helpers.HashPassword(req.Password)
	if err != nil {
		log.Printf("Internal Error: %v", err)
		return constants.ErrInternalErr
	}

	user := entities.User{
		Email:    req.Email,
		Name:     req.Name,
		Password: hashedPassword,
		RoleID:   role.ID,
	}

	if _, err := s.userRepo.Create(ctx, s.db, user); err != nil {
		log.Printf("Internal Error: %v", err)
		return constants.ErrInternalErr
	}
	return nil
}

func (s *userService) UpdateUser(ctx context.Context, id uuid.UUID, req dto.UserUpdateRequest) (dto.UserResponse, error) {
	user, err := s.userRepo.FindByID(ctx, s.db, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return dto.UserResponse{}, dto.ErrUserNotFound
		}
		log.Printf("Internal Error: %v", err)
		return dto.UserResponse{}, constants.ErrInternalErr
	}

	data := entities.User{}

	if req.Name != "" {
		data.Name = req.Name
	}
	if req.Email != "" {
		data.Email = req.Email
	}
	if req.Password != "" {
		if len(req.Password) < 8 {
			return dto.UserResponse{}, errors.New("password minimal 8 karakter")
		}
		hashedPassword, err := helpers.HashPassword(req.Password)
		if err != nil {
			log.Printf("Internal Error: %v", err)
			return dto.UserResponse{}, constants.ErrInternalErr
		}
		data.Password = hashedPassword
	}
	if req.RoleName != "" {
		role, err := s.roleRepo.FindByName(ctx, s.db, req.RoleName)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return dto.UserResponse{}, errors.New("role not found")
			}
			log.Printf("Internal Error: %v", err)
			return dto.UserResponse{}, constants.ErrInternalErr
		}
		data.RoleID = role.ID
	}

	updated, err := s.userRepo.Update(ctx, s.db, user.ID, data)
	if err != nil {
		log.Printf("Internal Error: %v", err)
		return dto.UserResponse{}, constants.ErrInternalErr
	}
	return toResponse(updated), nil
}

func (s *userService) DeleteUser(ctx context.Context, id uuid.UUID) error {
	if err := s.userRepo.Delete(ctx, s.db, id); err != nil {
		log.Printf("Internal Error: %v", err)
		return constants.ErrInternalErr
	}
	return nil
}
