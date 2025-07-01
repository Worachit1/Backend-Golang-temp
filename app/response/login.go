package response

type FindByStudentNumber struct {
	StudentNumber string `bun:"student_number" json:"student_number"`
}

type LoginResponse struct {
	ID            string `bun:"id" json:"id"`
	FirstName     string `bun:"first_name" json:"first_name"`
	LastName      string `bun:"last_name" json:"last_name"`
	StudentNumber string `bun:"student_number" json:"student_number"`
}
