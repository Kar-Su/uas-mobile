package dto

import "errors"

const (
	MESSAGE_FAILED_GET_SATUAN_BARANG    = "failed to get satuan barang"
	MESSAGE_FAILED_CREATE_SATUAN_BARANG = "failed to create satuan barang"
	MESSAGE_FAILED_UPDATE_SATUAN_BARANG = "failed to update satuan barang"
	MESSAGE_FAILED_DELETE_SATUAN_BARANG = "failed to delete satuan barang"
	MESSAGE_FAILED_BAD_REQUEST          = "bad request"

	MESSAGE_SUCCESS_GET_SATUAN_BARANG    = "success get satuan barang"
	MESSAGE_SUCCESS_CREATE_SATUAN_BARANG = "success create satuan barang"
	MESSAGE_SUCCESS_UPDATE_SATUAN_BARANG = "success update satuan barang"
	MESSAGE_SUCCESS_DELETE_SATUAN_BARANG = "success delete satuan barang"
)

var (
	ErrSatuanBarangNotFound    = errors.New("satuan barang not found")
	ErrSatuanBarangInUse       = errors.New("satuan barang masih digunakan oleh barang, tidak dapat dihapus")
)

type (
	SatuanBarangRequest struct {
		Satuan     string `json:"satuan" binding:"required"`
		Keterangan string `json:"keterangan"`
	}

	SatuanBarangResponse struct {
		ID         int     `json:"id"`
		Satuan     *string `json:"satuan"`
		Keterangan *string `json:"keterangan"`
	}
)
