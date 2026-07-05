package utils

const DefaultPageSize = 10

type PaginationQuery struct {
	Page  int `form:"page,default=1" example:"1"`
	Limit int `form:"limit" example:"10"`
}

type PaginationMeta struct {
	Page       int   `json:"page" example:"1"`
	PerPage    int   `json:"per_page" example:"10"`
	TotalItems int64 `json:"total_items" example:"100"`
	TotalPages int   `json:"total_pages" example:"10"`
}

func ResolvePageSize(limit int) int {
	if limit > 0 {
		return limit
	}
	return DefaultPageSize
}

func GetOffset(page int) int {
	return GetOffsetWithSize(page, DefaultPageSize)
}

func GetOffsetWithSize(page, pageSize int) int {
	if page <= 0 {
		page = 1
	}
	return (page - 1) * pageSize
}

func GetPage(page int) int {
	if page <= 0 {
		return 1
	}
	return page
}

func CountTotalPages(totalItems int64) int {
	return CountTotalPagesWithSize(totalItems, DefaultPageSize)
}

func CountTotalPagesWithSize(totalItems int64, pageSize int) int {
	if totalItems == 0 || pageSize <= 0 {
		return 1
	}
	pages := int(totalItems) / pageSize
	if int(totalItems)%pageSize != 0 {
		pages++
	}
	return pages
}

// NewPaginationMeta builds a PaginationMeta from a page number and total item count.
// Pass the result to WithMeta() when calling BuildResponse.
func NewPaginationMeta(page int, totalItems int64) *PaginationMeta {
	return NewPaginationMetaWithSize(page, totalItems, DefaultPageSize)
}

func NewPaginationMetaWithSize(page int, totalItems int64, pageSize int) *PaginationMeta {
	resolved := ResolvePageSize(pageSize)
	return &PaginationMeta{
		Page:       GetPage(page),
		PerPage:    resolved,
		TotalItems: totalItems,
		TotalPages: CountTotalPagesWithSize(totalItems, resolved),
	}
}
