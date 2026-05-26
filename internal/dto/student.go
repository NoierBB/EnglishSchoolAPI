package dto

type RegisterStudentRequest struct {
	Email    string `json:"email"`
	Password string `json:"password_hash"`
	Name     string `json:"name"`
	Age      int    `json:"age"`
	Level    string `json:"level"`
}
