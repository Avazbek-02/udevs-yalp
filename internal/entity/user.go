package entity

type User struct {
	ID          string `json:"id"`
	FullName    string `json:"full_name"`
	Username    string `json:"username"`
	Email       string `json:"email"`
	Password    string `json:"password"`
	UserType    string `json:"user_type"`
	UserRole    string `json:"user_role"`
	Status      string `json:"status"`
	Gender      string `json:"gender"`
	Bio         string `json:"bio"`
	AvatarId    string `json:"profile_picture"`
	AccessToken string `json:"access_token"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}

type UserSingleRequest struct {
	ID       string `json:"id"`
	Email    string `json:"email"`
	UserName string `json:"user_name"`
	UserType string `json:"user_type"`
}

type UserList struct {
	Items []User `json:"users"`
	Count int    `json:"count"`
}
