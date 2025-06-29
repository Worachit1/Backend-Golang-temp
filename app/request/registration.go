package request

type CreateRegistration struct {
	ActivitiesID string `json:"activity_id"`
	StudentsID   string `json:"Student_id"`
}

type UpdateRegistration struct {
	CreateRegistration
}

type ListRegistration struct {
	Page     int    `form:"page"`
	Size     int    `form:"size"`
	Search   string `form:"search"`
	SearchBy string `form:"search_by"`
	SortBy   string `form:"sort_by"`
	OrderBy  string `form:"order_by"`
}

type GetByIDRegistration struct {
	ID string `uri:"id" binding:"required"`
}
