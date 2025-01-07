package entity

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type RegisterRequest struct {
	FullName string `json:"full_name"`
	Username string `json:"user_name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type VerifyPhoneRequest struct {
	Email string `json:"email"`
	Otp         string `json:"otp"`
}
