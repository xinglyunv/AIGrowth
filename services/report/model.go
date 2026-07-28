package report

import "time"

type Report struct {
	ID             string                 `json:"id"`
	ProjectID      string                 `json:"project_id"`
	TaskID         string                 `json:"task_id"`
	UserID         string                 `json:"user_id"`
	Title          string                 `json:"title"`
	Type           string                 `json:"type"`
	VisibilityScore int                   `json:"visibility_score"`
	Content        map[string]interface{} `json:"content"`
	Summary        string                 `json:"summary"`
	Status         string                 `json:"status"`
	ShareToken     string                 `json:"share_token,omitempty"`
	ShareExpiresAt *time.Time             `json:"share_expires_at,omitempty"`
	CreatedAt      time.Time              `json:"created_at"`
	UpdatedAt      time.Time              `json:"updated_at"`
}
