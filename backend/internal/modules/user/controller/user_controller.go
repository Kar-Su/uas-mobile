package controller

import (
	"errors"
	"net/http"

	"github.com/Kar-Su/uas-mobile.git/internal/modules/user/dto"
	"github.com/Kar-Su/uas-mobile.git/internal/modules/user/repository"
	"github.com/Kar-Su/uas-mobile.git/internal/modules/user/service"
	"github.com/Kar-Su/uas-mobile.git/internal/package/constants"
	"github.com/Kar-Su/uas-mobile.git/internal/package/sse"
	"github.com/Kar-Su/uas-mobile.git/internal/package/utils"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/samber/do/v2"
	"gorm.io/gorm"
)

type UserController interface {
	Me(ctx *gin.Context)
	GetUsers(ctx *gin.Context)
	GetUser(ctx *gin.Context)
	CreateUser(ctx *gin.Context)
	UpdateUser(ctx *gin.Context)
	DeleteUser(ctx *gin.Context)
}

type userController struct {
	userService service.UserService
	db          *gorm.DB
}

func NewUserController(injector do.Injector, db *gorm.DB, userService service.UserService, roleRepo repository.RoleRepository) UserController {
	return &userController{
		userService: userService,
		db:          db,
	}
}

func (c *userController) Me(ctx *gin.Context) {
	path := ctx.Request.URL.Path
	userIDStr := ctx.MustGet("user_id").(string)

	id, err := uuid.Parse(userIDStr)
	if err != nil {
		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_GET_USER, err.Error(), utils.WithPath(path))
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	result, err := c.userService.GetUserByID(ctx.Request.Context(), id)
	if err != nil {
		if errors.Is(err, constants.ErrInternalErr) {
			res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_GET_USER, err.Error(), utils.WithPath(path))
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, res)
			return
		}
		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_GET_USER, err.Error(), utils.WithPath(path))
		ctx.AbortWithStatusJSON(http.StatusNotFound, res)
		return
	}

	res := utils.BuildResponse(dto.MESSAGE_SUCCESS_GET_USER, result, utils.WithPath(path))
	ctx.JSON(http.StatusOK, res)
}

func (c *userController) GetUsers(ctx *gin.Context) {
	path := ctx.Request.URL.Path

	var pageQuery utils.PaginationQuery
	if err := ctx.ShouldBindQuery(&pageQuery); err != nil {
		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_BAD_REQUEST, err.Error(), utils.WithPath(path))
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	var searchFilterQuery dto.SearchFilterQuery
	if err := ctx.ShouldBindQuery(&searchFilterQuery); err != nil {
		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_BAD_REQUEST, err.Error(), utils.WithPath(path))
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	var result any
	var total int64
	var err error
	if utils.IsNilStruct(&searchFilterQuery) {
		result, total, err = c.userService.GetAllUsers(ctx.Request.Context(), pageQuery.Page, nil)
	} else {
		result, total, err = c.userService.GetAllUsers(ctx.Request.Context(), pageQuery.Page, &searchFilterQuery)
	}

	if err != nil {
		if errors.Is(err, constants.ErrInternalErr) {
			res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_GET_USER, err.Error(), utils.WithPath(path))
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, res)
			return
		}
		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_GET_USER, err.Error(), utils.WithPath(path))
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	res := utils.BuildResponse(dto.MESSAGE_SUCCESS_GET_USER, result,
		utils.WithPath(path),
		utils.WithMeta(utils.NewPaginationMeta(pageQuery.Page, total)),
	)
	ctx.JSON(http.StatusOK, res)
}

func (c *userController) GetUser(ctx *gin.Context) {
	path := ctx.Request.URL.Path

	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_GET_USER, err.Error(), utils.WithPath(path))
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	result, err := c.userService.GetUserByID(ctx.Request.Context(), id)
	if err != nil {
		if errors.Is(err, constants.ErrInternalErr) {
			res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_GET_USER, err.Error(), utils.WithPath(path))
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, res)
			return
		}
		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_GET_USER, err.Error(), utils.WithPath(path))
		ctx.AbortWithStatusJSON(http.StatusNotFound, res)
		return
	}

	res := utils.BuildResponse(dto.MESSAGE_SUCCESS_GET_USER, result, utils.WithPath(path))
	ctx.JSON(http.StatusOK, res)
}

func (c *userController) CreateUser(ctx *gin.Context) {
	path := ctx.Request.URL.Path

	var req dto.UserCreateRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_BAD_REQUEST, utils.TranslateValidationError(err), utils.WithPath(path))
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	if err := c.userService.CreateUser(ctx.Request.Context(), req); err != nil {
		if errors.Is(err, constants.ErrInternalErr) {
			res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_CREATE_USER, err.Error(), utils.WithPath(path))
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, res)
			return
		}
		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_CREATE_USER, err.Error(), utils.WithPath(path))
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	res := utils.BuildResponse(dto.MESSAGE_SUCCESS_CREATE_USER, any(nil), utils.WithPath(path))
	sse.Default().Broadcast("user")
	ctx.JSON(http.StatusCreated, res)
}

func (c *userController) UpdateUser(ctx *gin.Context) {
	path := ctx.Request.URL.Path

	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_UPDATE_USER, err.Error(), utils.WithPath(path))
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	var req dto.UserUpdateRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_BAD_REQUEST, utils.TranslateValidationError(err), utils.WithPath(path))
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	result, err := c.userService.UpdateUser(ctx.Request.Context(), id, req)
	if err != nil {
		if errors.Is(err, constants.ErrInternalErr) {
			res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_UPDATE_USER, err.Error(), utils.WithPath(path))
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, res)
			return
		}
		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_UPDATE_USER, err.Error(), utils.WithPath(path))
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	res := utils.BuildResponse(dto.MESSAGE_SUCCESS_UPDATE_USER, result, utils.WithPath(path))
	sse.Default().Broadcast("user")
	ctx.JSON(http.StatusOK, res)
}

func (c *userController) DeleteUser(ctx *gin.Context) {
	path := ctx.Request.URL.Path

	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_DELETE_USER, err.Error(), utils.WithPath(path))
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	if err := c.userService.DeleteUser(ctx.Request.Context(), id); err != nil {
		if errors.Is(err, constants.ErrInternalErr) {
			res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_DELETE_USER, err.Error(), utils.WithPath(path))
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, res)
			return
		}
		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_DELETE_USER, err.Error(), utils.WithPath(path))
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	res := utils.BuildResponse(dto.MESSAGE_SUCCESS_DELETE_USER, any(nil), utils.WithPath(path))
	sse.Default().Broadcast("user")
	ctx.JSON(http.StatusOK, res)
}
