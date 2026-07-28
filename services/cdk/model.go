package cdk

import "time"

type CDKCode struct {
	ID        string     `json:"id"`
	Code      string     `json:"code"`
	Credits   int        `json:"credits"`
	MaxUses   int        `json:"max_uses"`
	UsedCount int        `json:"used_count"`
	IsActive  bool       `json:"is_active"`
	ExpiresAt *time.Time `json:"expires_at,omitempty"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
}

type CreateCDKRequest struct {
	Code      string     `json:"code"`
	Credits   int        `json:"credits"`
	MaxUses   int        `json:"max_uses"`
	ExpiresAt *time.Time `json:"expires_at,omitempty"`
}

type UpdateCDKRequest struct {
	IsActive  *bool      `json:"is_active,omitempty"`
	MaxUses   *int       `json:"max_uses,omitempty"`
	ExpiresAt *time.Time `json:"expires_at,omitempty"`
}

type CDKUsage struct {
	ID     string    `json:"id"`
	CDKID  string    `json:"cdk_id"`
	UserID string    `json:"user_id"`
	UsedAt time.Time `json:"used_at"`
}

type RedeemRequest struct {
	Code string `json:"code"`
}

type RedeemResult struct {
	Credits int    `json:"credits"`
	Message string `json:"message"`
}
