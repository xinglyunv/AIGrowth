package space

import "time"

const (
	RoleOwner  = "owner"
	RoleAdmin  = "admin"
	RoleMember = "member"
	RoleViewer = "viewer"
)

type Space struct {
	ID        string    `json:"id"`
	OwnerID   string    `json:"owner_id"`
	Name      string    `json:"name"`
	Slug      string    `json:"slug"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Member struct {
	UserID    string    `json:"user_id"`
	Email     string    `json:"email"`
	Username  string    `json:"username"`
	Role      string    `json:"role"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
}

type CreateRequest struct {
	Name string `json:"name"`
}

type InviteRequest struct {
	Email string `json:"email"`
	Role  string `json:"role"`
}

func IsValidRole(role string, allowOwner bool) bool {
	if role == RoleOwner {
		return allowOwner
	}
	return role == RoleAdmin || role == RoleMember || role == RoleViewer
}
