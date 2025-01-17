package entity

type Event struct {
	ID          string `json:"id"`
	BusinessID  string `json:"business_id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Date        string `json:"date"`
	Location    string `json:"location"`
	CreatedAt   string `json:"created_at"`
}

type EventParticipant struct {
	ID       string `json:"id"`
	EventID  string `json:"event_id"`
	UserID   string `json:"user_id"`
	JoinedAt string `json:"joined_at"`
}

type EventList struct {
	Events []Event `json:"events"`
	Count  int     `json:"count"`
}

type EventParticipantList struct {
	Participants []EventParticipant `json:"participants"`
	Count        int                `json:"count"`
}