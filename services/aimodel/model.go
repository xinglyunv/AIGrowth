package aimodel

import "time"

type AIModel struct {
	ID             string     `json:"id"`
	Name           string     `json:"name"`
	Provider       string     `json:"provider"`
	Model          string     `json:"model"`
	BaseURL        string     `json:"base_url"`
	APIKey         string     `json:"api_key,omitempty"`
	Enabled        bool       `json:"enabled"`
	Description    string     `json:"description"`
	IsSystem       bool       `json:"is_system"`
	LastTestedAt   *time.Time `json:"last_tested_at"`
	LastTestStatus string     `json:"last_test_status"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
}

type CreateModelRequest struct {
	Name        string `json:"name"`
	Provider    string `json:"provider"`
	Model       string `json:"model"`
	BaseURL     string `json:"base_url"`
	APIKey      string `json:"api_key"`
	Enabled     bool   `json:"enabled"`
	Description string `json:"description"`
}

type UpdateModelRequest struct {
	Name        *string `json:"name"`
	Provider    *string `json:"provider"`
	Model       *string `json:"model"`
	BaseURL     *string `json:"base_url"`
	APIKey      *string `json:"api_key"`
	Enabled     *bool   `json:"enabled"`
	Description *string `json:"description"`
}

type DiscoverRequest struct {
	BaseURL string `json:"base_url"`
	APIKey  string `json:"api_key"`
	ModelID string `json:"model_id"`
}

type DiscoveredModel struct {
	ID       string `json:"id"`
	Object   string `json:"object"`
	OwnedBy  string `json:"owned_by"`
}
