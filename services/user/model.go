package user

import (
	"time"
)

// User maps to the users table.
type User struct {
	ID            string     `json:"id"`
	Email         string     `json:"email"`
	Phone         *string    `json:"phone,omitempty"`
	PasswordHash  string     `json:"-"`
	Username      string     `json:"username"`
	CompanyName   *string    `json:"company_name,omitempty"`
	AvatarURL     *string    `json:"avatar_url,omitempty"`
	Role          string     `json:"role"`
	EmailVerified bool       `json:"email_verified"`
	Status        string     `json:"status"`
	Credits       int        `json:"credits"`
	LastLoginAt   *time.Time `json:"last_login_at,omitempty"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
}

// CreateUserRequest is the request body for user registration.
type CreateUserRequest struct {
	Email       string  `json:"email"`
	Phone       *string `json:"phone,omitempty"`
	Password    string  `json:"password"`
	Username    string  `json:"username"`
	CompanyName *string `json:"company_name,omitempty"`
}

// LoginRequest is the request body for user login.
type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// UpdateUserRequest is the request body for profile updates.
type UpdateUserRequest struct {
	Username    *string `json:"username,omitempty"`
	CompanyName *string `json:"company_name,omitempty"`
	AvatarURL   *string `json:"avatar_url,omitempty"`
}
