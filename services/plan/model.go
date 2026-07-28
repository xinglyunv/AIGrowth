package plan

import "time"

type Plan struct {
	ID           string    `json:"id"`
	Name         string    `json:"name"`
	Code         string    `json:"code"`
	Description  string    `json:"description"`
	MonthlyPrice float64   `json:"monthly_price"`
	YearlyPrice  float64   `json:"yearly_price"`
	MaxProjects  int       `json:"max_projects"`
	MaxAIQueries int       `json:"max_ai_queries"`
	MaxReports   int       `json:"max_reports"`
	Credits      int       `json:"credits"`
	Features     string    `json:"features"`
	Popular      bool      `json:"popular"`
	IsActive     bool      `json:"is_active"`
	SortOrder    int       `json:"sort_order"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type CreatePlanRequest struct {
	Name         string  `json:"name"`
	Code         string  `json:"code"`
	Description  string  `json:"description,omitempty"`
	MonthlyPrice float64 `json:"monthly_price"`
	YearlyPrice  float64 `json:"yearly_price"`
	MaxProjects  int     `json:"max_projects"`
	MaxAIQueries int     `json:"max_ai_queries"`
	MaxReports   int     `json:"max_reports"`
	Credits      int     `json:"credits"`
	SortOrder    int     `json:"sort_order"`
}

type UpdatePlanRequest struct {
	Name         *string  `json:"name,omitempty"`
	Code         *string  `json:"code,omitempty"`
	Description  *string  `json:"description,omitempty"`
	MonthlyPrice *float64 `json:"monthly_price,omitempty"`
	YearlyPrice  *float64 `json:"yearly_price,omitempty"`
	MaxProjects  *int     `json:"max_projects,omitempty"`
	MaxAIQueries *int     `json:"max_ai_queries,omitempty"`
	MaxReports   *int     `json:"max_reports,omitempty"`
	Credits      *int     `json:"credits,omitempty"`
	Features     *string  `json:"features,omitempty"`
	Popular      *bool    `json:"popular,omitempty"`
	IsActive     *bool    `json:"is_active,omitempty"`
	SortOrder    *int     `json:"sort_order,omitempty"`
}
