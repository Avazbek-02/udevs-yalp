package entity

type Review struct {
	ID         string `json:"id"`
	UserID     string `json:"user_id"`
	BusinessID string `json:"business_id"`
	Rating     int    `json:"rating"` // Value between 1 and 5
	Feedback   string `json:"feedback"`
	Photos     string `json:"photos"` // Assuming JSONB data is stored as a string
	CreatedAt  string `json:"created_at"`
}

type ReviewList struct {
	Items []Review `json:"reviews"`
	Count int      `json:"count"`
}
