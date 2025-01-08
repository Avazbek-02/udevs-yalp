package entity

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Platform string `json:"platform"`
}

type RegisterRequest struct {
	FullName string `json:"full_name"`
	Username string `json:"user_name"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Gender   string `json:"gender"`
}


type VerifyEmail struct {
	Email    string `json:"email"`
	Otp      string `json:"otp"`
	Platform string `json:"platform"`
}
