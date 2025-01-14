package entity

type Report struct {
	ID         string `json:"id"`
	UserID     string `json:"user_id"`
	BusinessID string `json:"business_id"`
	Reason     string `json:"reason"`
	CreatedAt  string `json:"created_at"`
}

type ReportList struct {
	Reports []Report `json:"reports"`
	Count   int      `json:"count"`
}
