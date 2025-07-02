package response

type ListUser struct {
	ID          string `bun:"id" json:"id"`
	FirstName   string `bun:"first_name" json:"first_name"`
	LastName    string `bun:"last_name" json:"last_name"`
	Email       string `bun:"email" json:"email"`
	Phone	   string `bun:"phone" json:"phone"`
	StudentNumber string `bun:"student_number" json:"student_number,omitempty"` // optional สำหรับนักศึกษา
	Address     string `bun:"address" json:"address"`
	CreatedAt   int64  `bun:"created_at" json:"create_at"`
	UpdatedAt   int64  `bun:"updated_at" json:"update_at"`
}
