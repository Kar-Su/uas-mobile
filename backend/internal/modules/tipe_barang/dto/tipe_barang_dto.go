package dto

import "errors"

const (
	MESSAGE_FAILED_GET_TIPE_BARANG    = "failed to get tipe barang"
	MESSAGE_FAILED_CREATE_TIPE_BARANG = "failed to create tipe barang"
	MESSAGE_FAILED_UPDATE_TIPE_BARANG = "failed to update tipe barang"
	MESSAGE_FAILED_DELETE_TIPE_BARANG = "failed to delete tipe barang"
	MESSAGE_FAILED_BAD_REQUEST        = "bad request"

	MESSAGE_SUCCESS_GET_TIPE_BARANG    = "success get tipe barang"
	MESSAGE_SUCCESS_CREATE_TIPE_BARANG = "success create tipe barang"
	MESSAGE_SUCCESS_UPDATE_TIPE_BARANG = "success update tipe barang"
	MESSAGE_SUCCESS_DELETE_TIPE_BARANG = "success delete tipe barang"
)

var (
	ErrTipeBarangNotFound    = errors.New("tipe barang not found")
	ErrTipeBarangInUse       = errors.New("tipe barang masih digunakan oleh barang, tidak dapat dihapus")
)

type (
	TipeBarangRequest struct {
		Name string `json:"name" binding:"required"`
	}

	TipeBarangResponse struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	}
)
