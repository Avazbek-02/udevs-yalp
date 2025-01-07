package entity

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type RegisterRequest struct {
	FullName    string `json:"full_name"`
	PhoneNumber string `json:"phone_number"`
}

type VerifyPhoneRequest struct {
	PhoneNumber string `json:"phone_number"`
	Otp         string `json:"otp"`
}

type ClientLoginRequest struct {
	PhoneNumber string `json:"phone_number"`
}
