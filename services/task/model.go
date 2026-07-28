package task

import (
	"strings"
	"time"
)

type AITask struct {
	ID             string     `json:"id"`
	ProjectID      string     `json:"project_id"`
	UserID         string     `json:"user_id"`
	Model          string     `json:"model"`
	Status         string     `json:"status"`
	QuestionsCount int        `json:"questions_count"`
	CompletedCount int        `json:"completed_count"`
	Progress       int        `json:"progress"`
	TotalTokens    int        `json:"total_tokens,omitempty"`
	TotalCost      float64    `json:"total_cost,omitempty"`
	ErrorMessage   string     `json:"error_message,omitempty"`
	StartedAt      *time.Time `json:"started_at,omitempty"`
	FinishedAt     *time.Time `json:"finished_at,omitempty"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
	ProjectName    string     `json:"project_name,omitempty"`
	Username       string     `json:"username,omitempty"`
}

type AIAnswer struct {
	ID             string                 `json:"id"`
	TaskID         string                 `json:"task_id"`
	Question       string                 `json:"question"`
	Answer         string                 `json:"answer"`
	Model          string                 `json:"model"`
	BrandMentioned bool                   `json:"brand_mentioned"`
	Sentiment      string                 `json:"sentiment,omitempty"`
	RankPosition   *int                   `json:"rank_position,omitempty"`
	Analysis       map[string]interface{} `json:"analysis,omitempty"`
	TokensUsed     int                    `json:"tokens_used,omitempty"`
	Cost           float64                `json:"cost,omitempty"`
	CreatedAt      time.Time              `json:"created_at"`
}

type CreateTaskRequest struct {
	ProjectID string `json:"project_id"`
	Model     string `json:"model,omitempty"`
}

// NormalizeModels converts the UI's comma-separated model selection into a stable list.
func NormalizeModels(value string) []string {
	seen := make(map[string]struct{})
	models := make([]string, 0)
	for _, raw := range strings.Split(value, ",") {
		model := strings.TrimSpace(raw)
		if model == "" {
			continue
		}
		if _, exists := seen[model]; exists {
			continue
		}
		seen[model] = struct{}{}
		models = append(models, model)
	}
	return models
}

func CreditCost(models []string) int {
	return len(models)
}

type TaskReport struct {
	Task            AITask       `json:"task"`
	Project         ProjectBrief `json:"project"`
	Answers         []AIAnswer   `json:"answers"`
	TotalQuestions  int          `json:"total_questions"`
	BrandMentions   int          `json:"brand_mentions"`
	VisibilityScore int          `json:"visibility_score"`
	Recommendations []string     `json:"recommendations"`
}

type ProjectBrief struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Industry string `json:"industry"`
	Website  string `json:"website,omitempty"`
}

type ModelComparisonResult struct {
	Model   string     `json:"model"`
	Status  string     `json:"status"`
	Answers []AIAnswer `json:"answers"`
	Score   float64    `json:"score"`
}

type ComparisonReport struct {
	Task    TaskReport              `json:"task"`
	Results []ModelComparisonResult `json:"results"`
}
