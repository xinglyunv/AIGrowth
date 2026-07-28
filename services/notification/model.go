package notification

import "time"

type Notification struct {
	ID        string    `json:"id"`
	UserID    string    `json:"user_id"`
	Type      string    `json:"type"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	IsRead    bool      `json:"is_read"`
	RelatedID string    `json:"related_id,omitempty"`
	CreatedAt time.Time `json:"created_at"`
}
