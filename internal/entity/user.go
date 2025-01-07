package entity

type User struct {
	ID          string `json:"id"`
	FullName    string `json:"full_name"`
	Username    string `json:"username"`
	Password    string `json:"password"`
	PhoneNumber string `json:"phone_number"`
	UserType    string `json:"user_type"`
	Status      string `json:"status"`
	AccessToken string `json:"access_token"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}

type UserSingleRequest struct {
	ID          string `json:"id"`
	PhoneNumber string `json:"phone_number"`
	UserName    string `json:"user_name"`
	UserType    string `json:"user_type"`
}

type UserList struct {
	Items []User `json:"users"`
	Count int    `json:"count"`
}
