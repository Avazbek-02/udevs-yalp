package entity

type Notification struct {
	ID        string `json:"id" db:"id"` // UUID of the notification
	OwnerId   string `json:"owner_id" db:"owner_id"`
	UserID    string `json:"user_id" db:"user_id"` // UUID of the user who the notification belongs to
	Email     string `json:"email" db:"email"`
	OwnerRole string `json:"ownerrole" db:"ownerrole"`
	Message   string `json:"message" db:"message"`       // Message of the notification
	Status    string `json:"status" db:"status"`         // Status of the notification ('read' or 'unread')
	CreatedAt string `json:"created_at" db:"created_at"` // Timestamp of when the notification was created
}

type NotificationList struct {
	Notifications []Notification `json:"notifications"` // List of notifications
	TotalCount    int            `json:"total_count"`   // Total number of notifications
}