package controller

import (
	"errors"
	"net/http"

	"github.com/Kar-Su/uas-mobile.git/internal/modules/barang/dto"
	"github.com/Kar-Su/uas-mobile.git/internal/modules/barang/service"
	"github.com/Kar-Su/uas-mobile.git/internal/package/constants"
	"github.com/Kar-Su/uas-mobile.git/internal/package/sse"
	"github.com/Kar-Su/uas-mobile.git/internal/package/utils"

	"github.com/gin-gonic/gin"
	"github.com/samber/do/v2"
	"gorm.io/gorm"
)

type BarangController interface {
	GetAll(ctx *gin.Context)
	GetByKode(ctx *gin.Context)
	Create(ctx *gin.Context)
	Update(ctx *gin.Context)
	Delete(ctx *gin.Context)
}

type barangController struct {
	service service.BarangService
	db      *gorm.DB
}

func NewBarangController(injector do.Injector, db *gorm.DB, svc service.BarangService) BarangController {
	return &barangController{service: svc, db: db}
}

func (c *barangController) GetAll(ctx *gin.Context) {
	path := ctx.Request.URL.Path

	var pageQuery utils.PaginationQuery
	if err := ctx.ShouldBindQuery(&pageQuery); err != nil {
		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_BAD_REQUEST, err.Error(), utils.WithPath(path))
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	var searchFilterQuery dto.FilterBarangQuery
	if err := ctx.ShouldBindQuery(&searchFilterQuery); err != nil {
		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_BAD_REQUEST, err.Error(), utils.WithPath(path))
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	result, total, err := c.service.GetAll(ctx.Request.Context(), pageQuery.Page, pageQuery.Limit, &searchFilterQuery)

	if err != nil {
		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_GET_BARANG, err.Error(), utils.WithPath(path))
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, res)
		return
	}

	res := utils.BuildResponse(dto.MESSAGE_SUCCESS_GET_BARANG, result,
		utils.WithPath(path),
		utils.WithMeta(utils.NewPaginationMetaWithSize(pageQuery.Page, total, pageQuery.Limit)),
	)
	ctx.JSON(http.StatusOK, res)
}

func (c *barangController) GetByKode(ctx *gin.Context) {
	path := ctx.Request.URL.Path
	kode := ctx.Param("kode")

	result, err := c.service.GetByKode(ctx.Request.Context(), kode)
	if err != nil {
		if errors.Is(err, dto.ErrBarangNotFound) {
			res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_GET_BARANG, err.Error(), utils.WithPath(path))
			ctx.AbortWithStatusJSON(http.StatusNotFound, res)
			return
		}
		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_GET_BARANG, err.Error(), utils.WithPath(path))
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, res)
		return
	}

	res := utils.BuildResponse(dto.MESSAGE_SUCCESS_GET_BARANG, result, utils.WithPath(path))
	ctx.JSON(http.StatusOK, res)
}

func (c *barangController) Create(ctx *gin.Context) {
	path := ctx.Request.URL.Path
	var req dto.BarangCreateRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_BAD_REQUEST, utils.TranslateValidationError(err), utils.WithPath(path))
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	result, err := c.service.Create(ctx.Request.Context(), req)
	if err != nil {
		if errors.Is(err, constants.ErrInternalErr) {
			res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_CREATE_BARANG, err.Error(), utils.WithPath(path))
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, res)
			return
		}
		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_CREATE_BARANG, err.Error(), utils.WithPath(path))
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	res := utils.BuildResponse(dto.MESSAGE_SUCCESS_CREATE_BARANG, result, utils.WithPath(path))
	sse.Default().Broadcast("changed")
	ctx.JSON(http.StatusCreated, res)
}

func (c *barangController) Update(ctx *gin.Context) {
	path := ctx.Request.URL.Path
	kode := ctx.Param("kode")

	var req dto.BarangUpdateRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_BAD_REQUEST, utils.TranslateValidationError(err), utils.WithPath(path))
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	result, err := c.service.Update(ctx.Request.Context(), kode, req)
	if err != nil {
		if errors.Is(err, dto.ErrBarangNotFound) {
			res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_UPDATE_BARANG, err.Error(), utils.WithPath(path))
			ctx.AbortWithStatusJSON(http.StatusNotFound, res)
			return
		}
		if errors.Is(err, constants.ErrInternalErr) {
			res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_UPDATE_BARANG, err.Error(), utils.WithPath(path))
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, res)
			return
		}
		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_UPDATE_BARANG, err.Error(), utils.WithPath(path))
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	res := utils.BuildResponse(dto.MESSAGE_SUCCESS_UPDATE_BARANG, result, utils.WithPath(path))
	sse.Default().Broadcast("changed")
	ctx.JSON(http.StatusOK, res)
}

func (c *barangController) Delete(ctx *gin.Context) {
	path := ctx.Request.URL.Path
	kode := ctx.Param("kode")

	if err := c.service.Delete(ctx.Request.Context(), kode); err != nil {
		if errors.Is(err, dto.ErrBarangNotFound) {
			res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_DELETE_BARANG, err.Error(), utils.WithPath(path))
			ctx.AbortWithStatusJSON(http.StatusNotFound, res)
			return
		}
		if errors.Is(err, constants.ErrInternalErr) {
			res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_DELETE_BARANG, err.Error(), utils.WithPath(path))
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, res)
			return
		}
		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_DELETE_BARANG, err.Error(), utils.WithPath(path))
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	res := utils.BuildResponse(dto.MESSAGE_SUCCESS_DELETE_BARANG, any(nil), utils.WithPath(path))
	sse.Default().Broadcast("changed")
	ctx.JSON(http.StatusOK, res)
}
