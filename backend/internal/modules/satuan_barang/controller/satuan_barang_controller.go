package controller

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/Kar-Su/uas-mobile.git/internal/modules/satuan_barang/dto"
	"github.com/Kar-Su/uas-mobile.git/internal/modules/satuan_barang/service"
	"github.com/Kar-Su/uas-mobile.git/internal/package/constants"
	"github.com/Kar-Su/uas-mobile.git/internal/package/sse"
	"github.com/Kar-Su/uas-mobile.git/internal/package/utils"

	"github.com/gin-gonic/gin"
	"github.com/samber/do/v2"
	"gorm.io/gorm"
)

type SatuanBarangController interface {
	GetAll(ctx *gin.Context)
	GetByID(ctx *gin.Context)
	Create(ctx *gin.Context)
	Update(ctx *gin.Context)
	Delete(ctx *gin.Context)
}

type satuanBarangController struct {
	service service.SatuanBarangService
	db      *gorm.DB
}

func NewSatuanBarangController(injector do.Injector, db *gorm.DB, svc service.SatuanBarangService) SatuanBarangController {
	return &satuanBarangController{service: svc, db: db}
}

func (c *satuanBarangController) GetAll(ctx *gin.Context) {
	path := ctx.Request.URL.Path
	result, err := c.service.GetAll(ctx.Request.Context())
	if err != nil {
		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_GET_SATUAN_BARANG, err.Error(), utils.WithPath(path))
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, res)
		return
	}
	res := utils.BuildResponse(dto.MESSAGE_SUCCESS_GET_SATUAN_BARANG, result, utils.WithPath(path))
	ctx.JSON(http.StatusOK, res)
}

func (c *satuanBarangController) GetByID(ctx *gin.Context) {
	path := ctx.Request.URL.Path
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_BAD_REQUEST, err.Error(), utils.WithPath(path))
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}
	result, err := c.service.GetByID(ctx.Request.Context(), id)
	if err != nil {
		if errors.Is(err, dto.ErrSatuanBarangNotFound) {
			res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_GET_SATUAN_BARANG, err.Error(), utils.WithPath(path))
			ctx.AbortWithStatusJSON(http.StatusNotFound, res)
			return
		}
		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_GET_SATUAN_BARANG, err.Error(), utils.WithPath(path))
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, res)
		return
	}
	res := utils.BuildResponse(dto.MESSAGE_SUCCESS_GET_SATUAN_BARANG, result, utils.WithPath(path))
	ctx.JSON(http.StatusOK, res)
}

func (c *satuanBarangController) Create(ctx *gin.Context) {
	path := ctx.Request.URL.Path
	var req dto.SatuanBarangRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_BAD_REQUEST, utils.TranslateValidationError(err), utils.WithPath(path))
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}
	result, err := c.service.Create(ctx.Request.Context(), req)
	if err != nil {
		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_CREATE_SATUAN_BARANG, err.Error(), utils.WithPath(path))
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, res)
		return
	}
	res := utils.BuildResponse(dto.MESSAGE_SUCCESS_CREATE_SATUAN_BARANG, result, utils.WithPath(path))
	sse.Default().Broadcast("satuan_barang")
	ctx.JSON(http.StatusCreated, res)
}

func (c *satuanBarangController) Update(ctx *gin.Context) {
	path := ctx.Request.URL.Path
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_BAD_REQUEST, err.Error(), utils.WithPath(path))
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}
	var req dto.SatuanBarangRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_BAD_REQUEST, utils.TranslateValidationError(err), utils.WithPath(path))
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}
	result, err := c.service.Update(ctx.Request.Context(), id, req)
	if err != nil {
		if errors.Is(err, dto.ErrSatuanBarangNotFound) {
			res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_UPDATE_SATUAN_BARANG, err.Error(), utils.WithPath(path))
			ctx.AbortWithStatusJSON(http.StatusNotFound, res)
			return
		}
		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_UPDATE_SATUAN_BARANG, err.Error(), utils.WithPath(path))
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, res)
		return
	}
	res := utils.BuildResponse(dto.MESSAGE_SUCCESS_UPDATE_SATUAN_BARANG, result, utils.WithPath(path))
	sse.Default().Broadcast("satuan_barang")
	ctx.JSON(http.StatusOK, res)
}

func (c *satuanBarangController) Delete(ctx *gin.Context) {
	path := ctx.Request.URL.Path
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_BAD_REQUEST, err.Error(), utils.WithPath(path))
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}
	if err := c.service.Delete(ctx.Request.Context(), id); err != nil {
		if errors.Is(err, dto.ErrSatuanBarangNotFound) {
			res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_DELETE_SATUAN_BARANG, err.Error(), utils.WithPath(path))
			ctx.AbortWithStatusJSON(http.StatusNotFound, res)
			return
		}
		if errors.Is(err, dto.ErrSatuanBarangInUse) {
			res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_DELETE_SATUAN_BARANG, err.Error(), utils.WithPath(path))
			ctx.AbortWithStatusJSON(http.StatusConflict, res)
			return
		}
		if errors.Is(err, constants.ErrInternalErr) {
			res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_DELETE_SATUAN_BARANG, err.Error(), utils.WithPath(path))
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, res)
			return
		}
		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_DELETE_SATUAN_BARANG, err.Error(), utils.WithPath(path))
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}
	res := utils.BuildResponse(dto.MESSAGE_SUCCESS_DELETE_SATUAN_BARANG, any(nil), utils.WithPath(path))
	sse.Default().Broadcast("satuan_barang")
	ctx.JSON(http.StatusOK, res)
}
