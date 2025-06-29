package response

type ListRegistration struct {
	ID           string `bun:"id" json:"id"`
	ActivitiesID string `bun:"activity_id" json:"activity_id"`
	StudentsID   string `bun:"student_id" json:"student_id"`
	CreatedAt    int64  `bun:"created_at" json:"created_at"`
	UpdatedAt    int64  `bun:"updated_at" json:"updated_at"`
}
