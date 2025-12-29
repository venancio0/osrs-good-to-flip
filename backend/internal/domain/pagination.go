package domain

// PaginationParams represents pagination parameters
type PaginationParams struct {
	Page  int // 1-based page number
	Limit int // Items per page
}

// PaginatedResult represents a paginated result
type PaginatedResult[T any] struct {
	Data       []T `json:"data"`
	Total      int `json:"total"`
	Page       int `json:"page"`
	Limit      int `json:"limit"`
	TotalPages int `json:"total_pages"`
}

// NewPaginationParams creates pagination params from query values
// Defaults: page=1, limit=20
func NewPaginationParams(page, limit int) PaginationParams {
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 20
	}
	if limit > 100 {
		limit = 100 // Max limit to prevent abuse
	}
	return PaginationParams{
		Page:  page,
		Limit: limit,
	}
}

// Offset calculates the offset for database queries
func (p PaginationParams) Offset() int {
	return (p.Page - 1) * p.Limit
}
