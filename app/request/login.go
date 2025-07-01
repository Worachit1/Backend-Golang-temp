package request

type LoginRequest struct {
	StudentNumber    string `json:"student_number"`
	Password string `json:"password"`
}