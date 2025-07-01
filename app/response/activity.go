package response

type ListActivity struct {
	ID          string `bun:"id" json:"id"`
	Name        string `bun:"name" json:"name"`
	Description string `bun:"description" json:"description"`
	ReleaseDate int64  `bun:"release_date" json:"release_date"`
	CreatedAt   int64  `bun:"created_at" json:"created_at"`
	UpdatedAt   int64  `bun:"updated_at" json:"updated_at"`
}
