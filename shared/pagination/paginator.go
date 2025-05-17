package pagination

import (
	"context"
	"strings"
)

type SortDirection string

const (
	ASC  SortDirection = "asc"
	DESC SortDirection = "desc"
)

type Sort struct {
	Field     string        `json:"field"`
	Direction SortDirection `json:"direction"`
}

func (s *Sort) ValidateAndSetDefaults() {
	if s.Direction == "" {
		s.Direction = ASC
	}
	if s.Direction != ASC && s.Direction != DESC {
		s.Direction = ASC
	}
}

type Pagination struct {
	Page     int    `json:"page"`
	PageSize int    `json:"size"`
	Total    int64  `json:"total"`
	Offset   int    `json:"-"`
	Sort     []Sort `json:"sort,omitempty"`
}
type PaginatedResponse struct {
	Data       any        `json:"data"`
	Pagination Pagination `json:"pagination,omitempty"`
}
type paginationKey struct{}

func WithPagination(ctx context.Context, p *Pagination) context.Context {
	return context.WithValue(ctx, paginationKey{}, p)
}

func GetPagination(ctx context.Context) (*Pagination, bool) {
	p, ok := ctx.Value(paginationKey{}).(*Pagination)
	return p, ok
}

func (p *Pagination) CalculateOffset() {
	p.Offset = (p.Page - 1) * p.PageSize
}
func SetPagination(ctx context.Context, page, size int, field, dir string) context.Context {
	var sorts []Sort
	if field != "" {
		sorts = append(sorts, Sort{
			Field:     field,
			Direction: SortDirection(strings.ToLower(dir)),
		})
	}

	
	defaultPage := 1
	defaultPageSize := 10

	// Check if page is valid, else use default
	if page <= 0 {
		page = defaultPage
	}

	// Check if size is valid, else use default
	if size <= 0 {
		size = defaultPageSize
	}

	pagination := &Pagination{
		Page:     page,
		PageSize: size,
		Offset:   (page - 1) * size,
		Sort:     sorts,
	}

	return WithPagination(ctx, pagination)
}
