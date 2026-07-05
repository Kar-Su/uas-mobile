package dto

import (
	"errors"

	"github.com/Kar-Su/uas-mobile.git/internal/package/utils"
	"github.com/google/uuid"
)

const (
	MESSAGE_FAILED_GET_USER    = "failed to get user"
	MESSAGE_FAILED_CREATE_USER = "failed to create user"
	MESSAGE_FAILED_UPDATE_USER = "failed to update user"
	MESSAGE_FAILED_DELETE_USER = "failed to delete user"
	MESSAGE_FAILED_LOGIN_USER  = "failed to login"
	MESSAGE_FAILED_BAD_REQUEST = "bad request"

	MESSAGE_SUCCESS_GET_USER    = "success get user"
	MESSAGE_SUCCESS_CREATE_USER = "success create user"
	MESSAGE_SUCCESS_UPDATE_USER = "success update user"
	MESSAGE_SUCCESS_DELETE_USER = "success delete user"
	MESSAGE_SUCCESS_LOGIN_USER  = "success login"
)

var ErrUserNotFound = errors.New("user not found")

type (
	UserLoginRequest struct {
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required"`
	}

	UserCreateRequest struct {
		Email    string `json:"email" binding:"required,email"`
		Name     string `json:"name" binding:"required"`
		Password string `json:"password" binding:"required,min=8"`
		RoleName string `json:"role_name" binding:"required"`
	}

	UserUpdateRequest struct {
		Name     string `json:"name"`
		Email    string `json:"email" binding:"omitempty,email"`
		Password string `json:"password" binding:"omitempty,min=8"`
		RoleName string `json:"role_name"`
	}

	UserResponse struct {
		ID       uuid.UUID `json:"id"`
		Email    string    `json:"email"`
		Name     string    `json:"name"`
		RoleName string    `json:"role_name"`
	}

	SearchFilterQuery struct {
		SearchQuery utils.SearchQuery `binding:"omitempty"`
		FilterRole  *string           `form:"role" binding:"omitempty"`
	}
)
