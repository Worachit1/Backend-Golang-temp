package request

import "time"

type CreateActivity struct {
	Name        string    `json:"name"`
	Description string    `json:"description"`
	ReleaseDate time.Time `json:"release_date"`
}

type UpdateActivity struct {
	CreateActivity
}

type ListActivity struct {
	Page     int    `form:"page"`
	Size     int    `form:"size"`
	Search   string `form:"search"`
	SearchBy string `form:"search_by"`
	SortBy   string `form:"sort_by"`
	OrderBy  string `form:"order_by"`
}

type GetByIDActivity struct {
	ID string `uri:"id" binding:"required"`
}
