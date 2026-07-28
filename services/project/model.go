package project

import "time"

// BrandProject maps to the brand_projects table, with optional brand_info joined.
type BrandProject struct {
	ID          string    `json:"id"`
	UserID      string    `json:"user_id"`
	Name        string    `json:"name"`
	Website     string    `json:"website,omitempty"`
	Industry    string    `json:"industry"`
	Description string    `json:"description,omitempty"`
	Keywords    []string  `json:"keywords,omitempty"`
	ServiceArea string    `json:"service_area,omitempty"`
	TargetUsers string    `json:"target_users,omitempty"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	BrandIntro  *string   `json:"brand_intro,omitempty"`
	ProductIntro *string  `json:"product_intro,omitempty"`
	ServiceIntro *string  `json:"service_intro,omitempty"`
	FAQ         *string   `json:"faq,omitempty"`
	Advantages  *string   `json:"advantages,omitempty"`
	Cases       *string   `json:"cases,omitempty"`
}

// CreateProjectRequest is the request body for creating a project.
type CreateProjectRequest struct {
	Name        string   `json:"name"`
	Website     string   `json:"website,omitempty"`
	Industry    string   `json:"industry"`
	Description string   `json:"description,omitempty"`
	Keywords    []string `json:"keywords,omitempty"`
	ServiceArea string   `json:"service_area,omitempty"`
	TargetUsers string   `json:"target_users,omitempty"`
}

// UpdateProjectRequest is the request body for updating a project.
type UpdateProjectRequest struct {
	Name        *string  `json:"name,omitempty"`
	Website     *string  `json:"website,omitempty"`
	Industry    *string  `json:"industry,omitempty"`
	Description *string  `json:"description,omitempty"`
	Keywords    []string `json:"keywords,omitempty"`
	ServiceArea *string  `json:"service_area,omitempty"`
	TargetUsers *string  `json:"target_users,omitempty"`
	BrandIntro  *string  `json:"brand_intro,omitempty"`
	ProductIntro *string `json:"product_intro,omitempty"`
	ServiceIntro *string `json:"service_intro,omitempty"`
	FAQ         *string  `json:"faq,omitempty"`
	Advantages  *string  `json:"advantages,omitempty"`
	Cases       *string  `json:"cases,omitempty"`
}
