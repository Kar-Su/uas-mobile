package dto

import (
	"errors"
	"time"

	"github.com/Kar-Su/uas-mobile.git/internal/package/utils"
)

const (
	MESSAGE_FAILED_GET_BARANG    = "failed to get barang"
	MESSAGE_FAILED_CREATE_BARANG = "failed to create barang"
	MESSAGE_FAILED_UPDATE_BARANG = "failed to update barang"
	MESSAGE_FAILED_DELETE_BARANG = "failed to delete barang"
	MESSAGE_FAILED_BAD_REQUEST   = "bad request"

	MESSAGE_SUCCESS_GET_BARANG    = "success get barang"
	MESSAGE_SUCCESS_CREATE_BARANG = "success create barang"
	MESSAGE_SUCCESS_UPDATE_BARANG = "success update barang"
	MESSAGE_SUCCESS_DELETE_BARANG = "success delete barang"
)

var ErrBarangNotFound = errors.New("barang not found")

type (
	BarangCreateRequest struct {
		Kode     string `json:"kode" binding:"required"`
		Name     string `json:"name" binding:"required"`
		TipeID   int    `json:"tipe_id" binding:"required"`
		SatuanID int    `json:"satuan_id" binding:"required"`
		Quantity int    `json:"quantity"`
	}

	BarangUpdateRequest struct {
		Name     string `json:"name"`
		TipeID   int    `json:"tipe_id"`
		SatuanID int    `json:"satuan_id"`
		Quantity *int   `json:"quantity"`
	}

	TipeResponse struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	}

	SatuanResponse struct {
		ID         int     `json:"id"`
		Satuan     *string `json:"satuan"`
		Keterangan *string `json:"keterangan"`
	}

	BarangResponse struct {
		Kode      string         `json:"kode"`
		Name      string         `json:"name"`
		Tipe      TipeResponse   `json:"tipe"`
		Satuan    SatuanResponse `json:"satuan"`
		Quantity  int            `json:"quantity"`
		CreatedAt time.Time      `json:"created_at"`
		UpdatedAt time.Time      `json:"updated_at"`
	}

	FilterBarangQuery struct {
		Search   *utils.SearchQuery `binding:"omitempty"`
		Tipe     *string            `form:"tipe" binding:"omitempty"`
		QtyOrder *string            `form:"qty_order" binding:"omitempty"`
	}
)
