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

type EventUsers struct {
	ID       string `json:"id"`
	EventID  string `json:"event_id"`
	FullName string `json:"full_name"`
	Username string `json:"username"`
	Email    string `json:"email"`
	UserType string `json:"user_type"`
	UserRole string `json:"user_role"`
	Status   string `json:"status"`
	Gender   string `json:"gender"`
}

type EventParticipantList struct {
	Participants []EventUsers `json:"participants"`
	Count        int          `json:"count"`
}
