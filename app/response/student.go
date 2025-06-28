package response

type ListStudent struct {
	ID             string `bun:"id" json:"id"`
	First_Name     string `bun:"first_name" json:"first_name"`
	Last_Name      string `bun:"last_name" json:"last_name"`
	Student_number string `bun:"student_number" json:"student_number"`
	CreatedAt      int64  `bun:"created_at" json:"created_at"`
	UpdatedAt      int64  `bun:"updated_at" json:"updated_at"`
}
