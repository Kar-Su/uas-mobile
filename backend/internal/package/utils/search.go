package utils

import (
	"gorm.io/gorm"
)

type SearchQuery struct {
	Search *string `form:"search" binding:"omitempty"`
}

func (s *SearchQuery) BindingQuery(query *gorm.DB, searchField []string) {
	if len(searchField) == 0 {
		return
	}
	queryString := ""
	args := make([]any, 0, len(searchField))
	keyword := "%" + *s.Search + "%"
	for i, field := range searchField {
		args = append(args, keyword)
		if i == 0 {
			queryString = field
		} else {
			queryString += " OR " + field
		}
		if i < len(searchField) {
			queryString += " ILIKE ?"
		}
	}
	query = query.Where(queryString, args...)
}
