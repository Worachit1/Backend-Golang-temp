package request

type CreateUser struct {
	FirstName     string `json:"first_name" form:"first_name" binding:"required"`
	LastName      string `json:"last_name" form:"last_name" binding:"required"`
	StudentNumber string `json:"student_number" form:"student_number"`
	Email         string `json:"email" form:"email" binding:"required,email"`
	Phone         string `json:"phone" form:"phone"`
	Address       string `json:"address" form:"address" binding:"required"`
	Password      string `json:"password" form:"password"`
}

type UpdateUser struct {
	CreateUser
}

type ListUser struct {
	Page     int    `form:"page"`
	Size     int    `form:"size"`
	Search   string `form:"search"`
	SearchBy string `form:"search_by"`
	SortBy   string `form:"sort_by"`
	OrderBy  string `form:"order_by"`
}

type GetByIdUser struct {
	ID string `uri:"id" binding:"required"`
}
