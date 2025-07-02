package response

type ListEmergency struct {
	ID          string  `json:"id"`
	UserID      string  `json:"user_id"`
	OfficerID   *string `json:"officer_id,omitempty"` // เจ้าหน้าที่ผู้รับเรื่อง
	Type        string  `json:"type"`                 // ประเภทเหตุฉุกเฉิน
	Status      string  `json:"status"`               // สถานะ
	Title       string  `json:"title"`
	Description string  `json:"description"`
	Location    string  `json:"location"`
	MapLink     string  `json:"map_link,omitempty"`
	FileURL     string  `json:"file_url,omitempty"`
	ActionNote  string  `json:"action_note,omitempty"`
	CreatedAt   int64   `json:"created_at"`
	UpdatedAt   int64   `json:"updated_at"`
}

type DetailEmergency struct {
	ID          string  `json:"id"`
	UserID      string  `json:"user_id"`
	OfficerID   *string `json:"officer_id,omitempty"` // เจ้าหน้าที่ผู้รับเรื่อง
	Type        string  `json:"type"`                 // ประเภทเหตุฉุกเฉิน
	Status      string  `json:"status"`               // สถานะ
	Title       string  `json:"title"`
	Description string  `json:"description"`
	Location    string  `json:"location"`
	MapLink     string  `json:"map_link,omitempty"`
	FileURL     string  `json:"file_url,omitempty"`
	ActionNote  string  `json:"action_note,omitempty"`
	CreatedAt   int64   `json:"created_at"`
	UpdatedAt   int64   `json:"updated_at"`
}

// Helper structs for user info
type UserInfo struct {
	ID        int64  `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	Phone     string `json:"phone"`
}
