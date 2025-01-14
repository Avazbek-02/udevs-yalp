package entity

type Business struct {
	ID          string `json:"id"`
	OwnerID     string `json:"owner_id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Category    string `json:"category"`
	Address     string `json:"address"`
	ContactInfo string `json:"contact_info"`
	Photos      string `json:"photos"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}

type BusinessList struct {
	Items []Business `json:"businesses"`
	Count int        `json:"count"`
}
