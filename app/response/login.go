package response

type FindByEmail struct {
	Email string `bun:"student_number" json:"email"`
}

type LoginResponse struct {
	ID            string `bun:"id" json:"id"`
	FirstName     string `bun:"first_name" json:"first_name"`
	LastName      string `bun:"last_name" json:"last_name"`
	StudentNumber string `bun:"student_number" json:"student_number,omitempty"` // optional สำหรับนักศึกษา
	Email         string `bun:"email" json:"email"`
	Phone         string `bun:"phone" json:"phone,omitempty"`
	Address       string `bun:"address" json:"address,omitempty"`
	Role          string `bun:"role" json:"role,omitempty"` // user, officer, admin
}
