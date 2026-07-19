package common

// PaginationQuery represents standard pagination parameters in request queries.
type PaginationQuery struct {
	Page    int    `form:"page" binding:"omitempty,min=1"`
	PerPage int    `form:"per_page" binding:"omitempty,min=1,max=100"`
	SortBy  string `form:"sort_by" binding:"omitempty"`
	Order   string `form:"order" binding:"omitempty,oneof=asc desc"`
}

// GetLimitOffset calculates limit and offset for SQL queries based on pagination parameters.
// If Page or PerPage are 0, it uses defaults (e.g. page 1, per_page 20).
func (p *PaginationQuery) GetLimitOffset() (limit, offset int) {
	page := p.Page
	if page < 1 {
		page = 1
	}

	limit = p.PerPage
	if limit < 1 {
		limit = 20
	} else if limit > 100 {
		limit = 100
	}

	offset = (page - 1) * limit
	return limit, offset
}

// GetOrder returns the order string for SQL queries.
func (p *PaginationQuery) GetOrder() string {
	if p.SortBy == "" {
		return "created_at desc" // Default fallback
	}

	order := p.Order
	if order != "asc" && order != "desc" {
		order = "asc"
	}
	
	// WARNING: When using this in a query, ensure SortBy is validated 
	// against a whitelist of columns to prevent SQL injection.
	return p.SortBy + " " + order
}
