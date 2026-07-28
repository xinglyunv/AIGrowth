package competitor

import "time"

type Competitor struct {
	ID           string                 `json:"id"`
	ProjectID    string                 `json:"project_id"`
	Name         string                 `json:"name"`
	Website      string                 `json:"website,omitempty"`
	MentionCount int                    `json:"mention_count"`
	RankPosition int                    `json:"rank_position"`
	Advantages   string                 `json:"advantages,omitempty"`
	Analysis     map[string]interface{} `json:"analysis,omitempty"`
	CreatedAt    time.Time              `json:"created_at"`
	UpdatedAt    time.Time              `json:"updated_at"`
}
