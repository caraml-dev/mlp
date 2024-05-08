package pagination

import (
	"fmt"
	"math"
)

// Paging can be used to capture paging information in API responses.
type Paging struct {
	Page  int32
	Pages int32
	Total int32
}

// PaginationOptions can be used to supply pagination filter options to APIs.
type PaginationOptions struct {
	Page     *int32 `json:"page,omitempty"`
	PageSize *int32 `json:"page_size,omitempty"`
}

// Paginator handles common pagination workflows.
type Paginator struct {
	DefaultPage     int32
	DefaultPageSize int32
	MaxPageSize     int32
}

func NewPaginator(defaultPage int32, defaultPageSize int32, maxPageSize int32) Paginator {
	return Paginator{
		DefaultPage:     defaultPage,
		DefaultPageSize: defaultPageSize,
		MaxPageSize:     maxPageSize,
	}
}

func (p Paginator) NewPaginationOptions(page *int32, pageSize *int32) PaginationOptions {
	if page == nil {
		page = &p.DefaultPage
	}
	if pageSize == nil {
		pageSize = &p.DefaultPageSize
	}

	return PaginationOptions{
		Page:     page,
		PageSize: pageSize,
	}
}

func (p Paginator) ValidatePaginationParams(page *int32, pageSize *int32) error {
	if pageSize != nil && (*pageSize <= 0 || *pageSize > p.MaxPageSize) {
		return fmt.Errorf("page size must be within range (0 < page_size <= %d) or unset",
			p.MaxPageSize)
	}
	if page != nil && *page <= 0 {
		return fmt.Errorf("page must be > 0 or unset")
	}

	return nil
}

func ToPaging(opts PaginationOptions, count int) *Paging {
	return &Paging{
		Page:  *opts.Page,
		Pages: int32(math.Ceil(float64(count) / float64(*opts.PageSize))),
		Total: int32(count),
	}
}
