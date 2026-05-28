package model

import "time"

// Pagination holds normalized paging parameters used across list queries.
type Pagination struct {
	Page     int `json:"page"`
	PageSize int `json:"page_size"`
}

// Offset returns the SQL OFFSET for the current page.
func (p Pagination) Offset() int { return (p.Page - 1) * p.PageSize }

// Normalize clamps page/size to safe defaults and limits.
func (p *Pagination) Normalize() {
	if p.Page < 1 {
		p.Page = 1
	}
	if p.PageSize < 1 {
		p.PageSize = 20
	}
	if p.PageSize > 100 {
		p.PageSize = 100
	}
}

// PagedResult is the standard envelope for paginated list endpoints.
type PagedResult[T any] struct {
	Items      []T   `json:"items"`
	Page       int   `json:"page"`
	PageSize   int   `json:"page_size"`
	TotalItems int64 `json:"total_items"`
	TotalPages int64 `json:"total_pages"`
}

// NewPagedResult builds a paged envelope, computing total pages.
func NewPagedResult[T any](items []T, p Pagination, total int64) PagedResult[T] {
	if items == nil {
		items = []T{}
	}
	pages := total / int64(p.PageSize)
	if total%int64(p.PageSize) != 0 {
		pages++
	}
	return PagedResult[T]{
		Items:      items,
		Page:       p.Page,
		PageSize:   p.PageSize,
		TotalItems: total,
		TotalPages: pages,
	}
}

// Timestamps embeds standard audit columns.
type Timestamps struct {
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
