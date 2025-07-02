package request

type CreateEmergency struct {
	Type        string `json:"type"`        // เช่น "ไฟไหม้", "อุบัติเหตุ"
	Title       string `json:"title"`
	Description string `json:"description"`
	Location    string `json:"location"`
	MapLink     string `json:"map_link"`
}

type UpdateEmergency struct {
	Type        *string `json:"type,omitempty" form:"type"`
	Status      *string `json:"status,omitempty" form:"status"` // เช่น "รอการตอบสนอง", "กำลังดำเนินการ", "เสร็จสิ้น"
	Title       *string `json:"title,omitempty" form:"title"`
	Description *string `json:"description,omitempty" form:"description"`
	Location    *string `json:"location,omitempty" form:"location"`
	MapLink     *string `json:"map_link,omitempty" form:"map_link"`
	FileURL     *string `json:"file_url,omitempty" form:"file_url"`
	ActionNote  *string `json:"action_note,omitempty" form:"action_note"` // บันทึกการดำเนินการ
	OfficerID   *string `json:"officer_id,omitempty" form:"officer_id"`   // เจ้าหน้าที่ผู้รับเรื่อง
}

type UpdateEmergencyByOfficer struct {
	Status     *string `json:"status,omitempty" form:"status"`         // เปลี่ยนสถานะ เช่น "กำลังดำเนินการ", "เสร็จสิ้น"
	ActionNote *string `json:"action_note,omitempty" form:"action_note"` // บันทึกการดำเนินการ
	OfficerID  *string `json:"officer_id,omitempty" form:"officer_id"`   // เจ้าหน้าที่ผู้รับเรื่อง
}

type ListEmergency struct {
	Page     int    `form:"page"`
	Size     int    `form:"size"`
	Search   string `form:"search"`
	SearchBy string `form:"search_by"`
	SortBy   string `form:"sort_by"`
	OrderBy  string `form:"order_by"`
	Type     string `form:"type"`    // กรองตามประเภท
	Status   string `form:"status"`  // กรองตามสถานะ
	UserID   string `form:"user_id"` // กรองตามผู้แจ้ง
}

type GetByIDEmergency struct {
	ID string `uri:"id" binding:"required"`
}
type GetByUserIDEmergency struct {
    UserID string `uri:"id" binding:"required"`
}

type CreateEmergencyRequest struct {
    UserID      string `json:"user_id" binding:"required"`
    Type        string `json:"type" binding:"required"`
    Title       string `json:"title" binding:"required"`
    Description string `json:"description" binding:"required"`
    Location    string `json:"location" binding:"required"`
    MapLink     string `json:"map_link"`
}