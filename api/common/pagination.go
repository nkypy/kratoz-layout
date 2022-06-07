package common

import (
	"math"

	"gorm.io/gorm"
)

const DefaultPerPage = 5

func NewPaginationReply(currentPage, totalCount, perPage int64) *PaginationReply {
	return &PaginationReply{
		CurrentPage: currentPage,
		TotalCount:  totalCount,
		PerPage:     perPage,
		TotalPages:  int64(int(math.Ceil(float64(totalCount) / float64(perPage)))),
	}
}

// Paginate paginate query
func (x *PaginationRequest) Paginate(db *gorm.DB) *PaginationReply {
	if x == nil {
		return nil
	}

	if x.Page <= 0 {
		x.Page = 1
	}

	if x.PerPage <= 0 {
		x.PerPage = DefaultPerPage
	}

	var totalCount int64 = 0
	db.Count(&totalCount)

	offset := x.PerPage * (x.Page - 1)
	db.Limit(int(x.PerPage)).Offset(int(offset))

	return NewPaginationReply(x.Page, totalCount, x.PerPage)
}
