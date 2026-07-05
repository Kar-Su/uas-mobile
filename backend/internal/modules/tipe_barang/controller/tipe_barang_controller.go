package controller

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/Kar-Su/uas-mobile.git/internal/modules/tipe_barang/dto"
	"github.com/Kar-Su/uas-mobile.git/internal/modules/tipe_barang/service"
	"github.com/Kar-Su/uas-mobile.git/internal/package/constants"
	"github.com/Kar-Su/uas-mobile.git/internal/package/sse"
	"github.com/Kar-Su/uas-mobile.git/internal/package/utils"

	"github.com/gin-gonic/gin"
	"github.com/samber/do/v2"
	"gorm.io/gorm"
)

type TipeBarangController interface {
	GetAll(ctx *gin.Context)
	GetByID(ctx *gin.Context)
	Create(ctx *gin.Context)
	Update(ctx *gin.Context)
	Delete(ctx *gin.Context)
}

type tipeBarangController struct {
	service service.TipeBarangService
	db      *gorm.DB
}

func NewTipeBarangController(injector do.Injector, db *gorm.DB, svc service.TipeBarangService) TipeBarangController {
	return &tipeBarangController{service: svc, db: db}
}

func (c *tipeBarangController) GetAll(ctx *gin.Context) {
	path := ctx.Request.URL.Path
	result, err := c.service.GetAll(ctx.Request.Context())
	if err != nil {
		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_GET_TIPE_BARANG, err.Error(), utils.WithPath(path))
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, res)
		return
	}
	res := utils.BuildResponse(dto.MESSAGE_SUCCESS_GET_TIPE_BARANG, result, utils.WithPath(path))
	ctx.JSON(http.StatusOK, res)
}

func (c *tipeBarangController) GetByID(ctx *gin.Context) {
	path := ctx.Request.URL.Path
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_BAD_REQUEST, err.Error(), utils.WithPath(path))
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}
	result, err := c.service.GetByID(ctx.Request.Context(), id)
	if err != nil {
		if errors.Is(err, dto.ErrTipeBarangNotFound) {
			res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_GET_TIPE_BARANG, err.Error(), utils.WithPath(path))
			ctx.AbortWithStatusJSON(http.StatusNotFound, res)
			return
		}
		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_GET_TIPE_BARANG, err.Error(), utils.WithPath(path))
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, res)
		return
	}
	res := utils.BuildResponse(dto.MESSAGE_SUCCESS_GET_TIPE_BARANG, result, utils.WithPath(path))
	ctx.JSON(http.StatusOK, res)
}

func (c *tipeBarangController) Create(ctx *gin.Context) {
	path := ctx.Request.URL.Path
	var req dto.TipeBarangRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_BAD_REQUEST, utils.TranslateValidationError(err), utils.WithPath(path))
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}
	result, err := c.service.Create(ctx.Request.Context(), req)
	if err != nil {
		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_CREATE_TIPE_BARANG, err.Error(), utils.WithPath(path))
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, res)
		return
	}
	res := utils.BuildResponse(dto.MESSAGE_SUCCESS_CREATE_TIPE_BARANG, result, utils.WithPath(path))
	sse.Default().Broadcast("tipe_barang")
	ctx.JSON(http.StatusCreated, res)
}

func (c *tipeBarangController) Update(ctx *gin.Context) {
	path := ctx.Request.URL.Path
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_BAD_REQUEST, err.Error(), utils.WithPath(path))
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}
	var req dto.TipeBarangRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_BAD_REQUEST, utils.TranslateValidationError(err), utils.WithPath(path))
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}
	result, err := c.service.Update(ctx.Request.Context(), id, req)
	if err != nil {
		if errors.Is(err, dto.ErrTipeBarangNotFound) {
			res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_UPDATE_TIPE_BARANG, err.Error(), utils.WithPath(path))
			ctx.AbortWithStatusJSON(http.StatusNotFound, res)
			return
		}
		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_UPDATE_TIPE_BARANG, err.Error(), utils.WithPath(path))
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, res)
		return
	}
	res := utils.BuildResponse(dto.MESSAGE_SUCCESS_UPDATE_TIPE_BARANG, result, utils.WithPath(path))
	sse.Default().Broadcast("tipe_barang")
	ctx.JSON(http.StatusOK, res)
}

func (c *tipeBarangController) Delete(ctx *gin.Context) {
	path := ctx.Request.URL.Path
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_BAD_REQUEST, err.Error(), utils.WithPath(path))
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}
	if err := c.service.Delete(ctx.Request.Context(), id); err != nil {
		if errors.Is(err, dto.ErrTipeBarangNotFound) {
			res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_DELETE_TIPE_BARANG, err.Error(), utils.WithPath(path))
			ctx.AbortWithStatusJSON(http.StatusNotFound, res)
			return
		}
		if errors.Is(err, dto.ErrTipeBarangInUse) {
			res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_DELETE_TIPE_BARANG, err.Error(), utils.WithPath(path))
			ctx.AbortWithStatusJSON(http.StatusConflict, res)
			return
		}
		if errors.Is(err, constants.ErrInternalErr) {
			res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_DELETE_TIPE_BARANG, err.Error(), utils.WithPath(path))
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, res)
			return
		}
		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_DELETE_TIPE_BARANG, err.Error(), utils.WithPath(path))
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}
	res := utils.BuildResponse(dto.MESSAGE_SUCCESS_DELETE_TIPE_BARANG, any(nil), utils.WithPath(path))
	sse.Default().Broadcast("tipe_barang")
	ctx.JSON(http.StatusOK, res)
}
