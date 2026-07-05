package repository

import (
	"context"
	"errors"

	"github.com/Kar-Su/uas-mobile.git/internal/database/entities"
	"github.com/Kar-Su/uas-mobile.git/internal/modules/user/dto"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserRepository interface {
	CheckEmail(ctx context.Context, tx *gorm.DB, email string) (entities.User, bool, error)
	FindAll(ctx context.Context, tx *gorm.DB, offset, limit int, searchFilterQuery *dto.SearchFilterQuery) ([]entities.User, int64, error)
	FindByID(ctx context.Context, tx *gorm.DB, id uuid.UUID) (entities.User, error)
	Create(ctx context.Context, tx *gorm.DB, user entities.User) (entities.User, error)
	Update(ctx context.Context, tx *gorm.DB, id uuid.UUID, data entities.User) (entities.User, error)
	Delete(ctx context.Context, tx *gorm.DB, id uuid.UUID) error
}

type userRepository struct {
	roleRepo RoleRepository
	db       *gorm.DB
}

func NewUserRepository(db *gorm.DB, roleRepo RoleRepository) UserRepository {
	return &userRepository{roleRepo: roleRepo, db: db}
}

func (r *userRepository) CheckEmail(ctx context.Context, tx *gorm.DB, email string) (entities.User, bool, error) {
	if tx == nil {
		tx = r.db
	}
	var user entities.User
	if err := tx.WithContext(ctx).Where("email = ?", email).Preload("Role").First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return entities.User{}, false, nil
		}
		return entities.User{}, false, err
	}
	return user, true, nil
}

func (r *userRepository) FindAll(ctx context.Context, tx *gorm.DB, offset, limit int, searchFilterQuery *dto.SearchFilterQuery) ([]entities.User, int64, error) {
	if tx == nil {
		tx = r.db
	}

	var users []entities.User

	query := tx.WithContext(ctx).Model(&entities.User{}).Preload("Role")
	if searchFilterQuery != nil {
		if searchFilterQuery.SearchQuery.Search != nil {
			searchFilterQuery.SearchQuery.BindingQuery(query, []string{"name", "email"})
		}
		if searchFilterQuery.FilterRole != nil {
			role, err := r.roleRepo.FindByName(ctx, tx, *searchFilterQuery.FilterRole)
			if err != nil {
				return nil, 0, err
			}
			query = query.Where("role_id = ?", role.ID)
		}
	}

	var total int64
	if err := query.Session(&gorm.Session{}).Select("id").Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := query.Offset(offset).Limit(limit).Find(&users).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, 0, nil
		}
		return nil, 0, err
	}

	return users, total, nil
}

func (r *userRepository) FindByID(ctx context.Context, tx *gorm.DB, id uuid.UUID) (entities.User, error) {
	if tx == nil {
		tx = r.db
	}
	var user entities.User
	if err := tx.WithContext(ctx).Preload("Role").First(&user, "id = ?", id).Error; err != nil {
		return entities.User{}, err
	}
	return user, nil
}

func (r *userRepository) Create(ctx context.Context, tx *gorm.DB, user entities.User) (entities.User, error) {
	if tx == nil {
		tx = r.db
	}
	if err := tx.WithContext(ctx).Create(&user).Error; err != nil {
		return entities.User{}, err
	}
	return r.FindByID(ctx, tx, user.ID)
}

func (r *userRepository) Update(ctx context.Context, tx *gorm.DB, id uuid.UUID, data entities.User) (entities.User, error) {
	if tx == nil {
		tx = r.db
	}
	if err := tx.WithContext(ctx).Model(&entities.User{}).Where("id = ?", id).Updates(&data).Error; err != nil {
		return entities.User{}, err
	}
	return r.FindByID(ctx, tx, id)
}

func (r *userRepository) Delete(ctx context.Context, tx *gorm.DB, id uuid.UUID) error {
	if tx == nil {
		tx = r.db
	}
	if err := tx.WithContext(ctx).Where("id = ?", id).Delete(&entities.User{}).Error; err != nil {
		return err
	}
	return nil
}
