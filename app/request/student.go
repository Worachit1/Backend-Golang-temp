package request

type CreateStudent struct {
	First_name     string `json:"first_name"`
	Last_name      string `json:"last_name"`
	Student_number string `json:"student_number"`
}

type UpdateStudent struct {
	CreateStudent
}

type ListStudent struct {
	Page     int    `form:"page"`
	Size     int    `form:"size"`
	Search   string `form:"search"`
	SearchBy string `form:"search_by"`
	SortBy   string `form:"sort_by"`
	OrderBy  string `form:"order_by"`
}

type GetByIDStudent struct {
	ID string `uri:"id" binding:"required"`
}
