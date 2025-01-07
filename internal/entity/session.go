package entity

type Session struct {
	ID           string `json:"id"`
	UserID       string `json:"user_id"`
	IPAddress    string `json:"ip_address"`
	UserAgent    string `json:"user_agent"`
	IsActive     bool   `json:"is_active"`
	ExpiresAt    string `json:"expires_at"`
	LastActiveAt string `json:"last_active_at"`
	CreatedAt    string `json:"created_at"`
	UpdatedAt    string `json:"updated_at"`
}

type SessionList struct {
	Items []Session `json:"sessions"`
	Count    int       `json:"count"`
}
