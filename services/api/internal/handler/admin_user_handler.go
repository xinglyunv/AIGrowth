package handler

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/aige/api/internal/middleware"
	"github.com/aige/audit"
	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"
)

type AdminUserHandler struct {
	Pool      *pgxpool.Pool
	AuditRepo audit.Repository
}

func (h *AdminUserHandler) ListUsers(w http.ResponseWriter, r *http.Request) {
	adminID := middleware.GetAdminID(r.Context())
	if adminID == "" {
		jsonResponse(w, http.StatusUnauthorized, map[string]interface{}{
			"success": false, "message": "unauthorized",
		})
		return
	}

	offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	if limit <= 0 || limit > 100 {
		limit = 20
	}
	if offset < 0 {
		offset = 0
	}

	var total int
	err := h.Pool.QueryRow(r.Context(), `SELECT COUNT(*) FROM users`).Scan(&total)
	if err != nil {
		jsonResponse(w, http.StatusInternalServerError, map[string]interface{}{
			"success": false, "message": "internal server error",
		})
		return
	}

	rows, err := h.Pool.Query(r.Context(),
		`SELECT id, email, COALESCE(phone, '') AS phone, username, COALESCE(company_name, '') AS company_name,
		        role, email_verified, status, last_login_at, created_at, updated_at, COALESCE(credits, 0)
		 FROM users ORDER BY created_at DESC LIMIT $1 OFFSET $2`,
		limit, offset,
	)
	if err != nil {
		jsonResponse(w, http.StatusInternalServerError, map[string]interface{}{
			"success": false, "message": "internal server error",
		})
		return
	}
	defer rows.Close()

	type UserBrief struct {
		ID            string     `json:"id"`
		Email         string     `json:"email"`
		Phone         string     `json:"phone"`
		Username      string     `json:"username"`
		CompanyName   string     `json:"company_name"`
		Role          string     `json:"role"`
		EmailVerified bool       `json:"email_verified"`
		Status        string     `json:"status"`
		LastLoginAt   *time.Time `json:"last_login_at"`
		CreatedAt     time.Time  `json:"created_at"`
		UpdatedAt     time.Time  `json:"updated_at"`
		Credits       int        `json:"credits"`
	}

	var users []UserBrief
	for rows.Next() {
		var u UserBrief
		err := rows.Scan(&u.ID, &u.Email, &u.Phone, &u.Username, &u.CompanyName,
			&u.Role, &u.EmailVerified, &u.Status, &u.LastLoginAt, &u.CreatedAt, &u.UpdatedAt, &u.Credits)
		if err != nil {
			jsonResponse(w, http.StatusInternalServerError, map[string]interface{}{
				"success": false, "message": "internal server error",
			})
			return
		}
		users = append(users, u)
	}
	if users == nil {
		users = []UserBrief{}
	}

	jsonResponse(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"data":    users,
		"meta": map[string]interface{}{
			"total":  total,
			"offset": offset,
			"limit":  limit,
		},
	})
}

func (h *AdminUserHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	adminID := middleware.GetAdminID(r.Context())
	if adminID == "" {
		jsonResponse(w, http.StatusUnauthorized, map[string]interface{}{
			"success": false, "message": "unauthorized",
		})
		return
	}

	id := chi.URLParam(r, "id")

	type UserDetail struct {
		ID            string     `json:"id"`
		Email         string     `json:"email"`
		Phone         string     `json:"phone"`
		Username      string     `json:"username"`
		CompanyName   string     `json:"company_name"`
		Role          string     `json:"role"`
		EmailVerified bool       `json:"email_verified"`
		Status        string     `json:"status"`
		LastLoginAt   *time.Time `json:"last_login_at"`
		CreatedAt     time.Time  `json:"created_at"`
		UpdatedAt     time.Time  `json:"updated_at"`
		Credits       int        `json:"credits"`
	}

	var u UserDetail
	err := h.Pool.QueryRow(r.Context(),
		`SELECT id, email, COALESCE(phone, '') AS phone, username, COALESCE(company_name, '') AS company_name,
			        role, email_verified, status, last_login_at, created_at, updated_at, COALESCE(credits, 0)
		 FROM users WHERE id = $1`, id,
	).Scan(&u.ID, &u.Email, &u.Phone, &u.Username, &u.CompanyName,
		&u.Role, &u.EmailVerified, &u.Status, &u.LastLoginAt, &u.CreatedAt, &u.UpdatedAt, &u.Credits)
	if err != nil {
		jsonResponse(w, http.StatusNotFound, map[string]interface{}{
			"success": false, "message": "user not found",
		})
		return
	}

	jsonResponse(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"data":    u,
	})
}

func (h *AdminUserHandler) GetUserProjects(w http.ResponseWriter, r *http.Request) {
	adminID := middleware.GetAdminID(r.Context())
	if adminID == "" {
		jsonResponse(w, http.StatusUnauthorized, map[string]interface{}{
			"success": false, "message": "unauthorized",
		})
		return
	}

	userID := chi.URLParam(r, "id")

	rows, err := h.Pool.Query(r.Context(),
		`SELECT id, user_id, name, COALESCE(website, '') AS website, industry,
		        COALESCE(description, '') AS description, keywords, status, created_at, updated_at
		 FROM brand_projects WHERE user_id = $1 ORDER BY created_at DESC`,
		userID,
	)
	if err != nil {
		jsonResponse(w, http.StatusInternalServerError, map[string]interface{}{
			"success": false, "message": "internal server error",
		})
		return
	}
	defer rows.Close()

	type ProjectBrief struct {
		ID          string   `json:"id"`
		UserID      string   `json:"user_id"`
		Name        string   `json:"name"`
		Website     string   `json:"website"`
		Industry    string   `json:"industry"`
		Description string   `json:"description"`
		Keywords    []string `json:"keywords"`
		Status      string   `json:"status"`
		CreatedAt   string   `json:"created_at"`
		UpdatedAt   string   `json:"updated_at"`
	}

	var projects []ProjectBrief
	for rows.Next() {
		var p ProjectBrief
		err := rows.Scan(&p.ID, &p.UserID, &p.Name, &p.Website, &p.Industry,
			&p.Description, &p.Keywords, &p.Status, &p.CreatedAt, &p.UpdatedAt)
		if err != nil {
			jsonResponse(w, http.StatusInternalServerError, map[string]interface{}{
				"success": false, "message": "internal server error",
			})
			return
		}
		projects = append(projects, p)
	}
	if projects == nil {
		projects = []ProjectBrief{}
	}

	jsonResponse(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"data":    projects,
	})
}

func (h *AdminUserHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	adminID := middleware.GetAdminID(r.Context())
	if adminID == "" {
		jsonResponse(w, http.StatusUnauthorized, map[string]interface{}{
			"success": false, "message": "unauthorized",
		})
		return
	}

	id := chi.URLParam(r, "id")

	var req struct {
		Username *string `json:"username"`
		Status   *string `json:"status"`
		Role     *string `json:"role"`
		Email    *string `json:"email"`
		Password *string `json:"password"`
		Credits  *int    `json:"credits"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		jsonResponse(w, http.StatusBadRequest, map[string]interface{}{
			"success": false, "message": "invalid request body",
		})
		return
	}

	updated := false

	if req.Username != nil {
		username := strings.TrimSpace(*req.Username)
		if username == "" {
			jsonResponse(w, http.StatusBadRequest, map[string]interface{}{
				"success": false, "message": "username cannot be empty",
			})
			return
		}
		_, err := h.Pool.Exec(r.Context(),
			`UPDATE users SET username = $2, updated_at = NOW() WHERE id = $1`,
			id, username,
		)
		if err != nil {
			jsonResponse(w, http.StatusInternalServerError, map[string]interface{}{
				"success": false, "message": "internal server error",
			})
			return
		}
		updated = true
	}

	if req.Status != nil {
		s := *req.Status
		if s != "active" && s != "disabled" {
			jsonResponse(w, http.StatusBadRequest, map[string]interface{}{
				"success": false, "message": "status must be 'active' or 'disabled'",
			})
			return
		}
		_, err := h.Pool.Exec(r.Context(),
			`UPDATE users SET status = $2, updated_at = NOW() WHERE id = $1`,
			id, s,
		)
		if err != nil {
			jsonResponse(w, http.StatusInternalServerError, map[string]interface{}{
				"success": false, "message": "internal server error",
			})
			return
		}
		updated = true
	}

	if req.Role != nil {
		r2 := *req.Role
		if r2 != "user" && r2 != "admin" && r2 != "superadmin" {
			jsonResponse(w, http.StatusBadRequest, map[string]interface{}{
				"success": false, "message": "role must be 'user', 'admin', or 'superadmin'",
			})
			return
		}
		_, err := h.Pool.Exec(r.Context(),
			`UPDATE users SET role = $2, updated_at = NOW() WHERE id = $1`,
			id, r2,
		)
		if err != nil {
			jsonResponse(w, http.StatusInternalServerError, map[string]interface{}{
				"success": false, "message": "internal server error",
			})
			return
		}
		updated = true
	}

	if req.Email != nil {
		_, err := h.Pool.Exec(r.Context(),
			`UPDATE users SET email = $2, updated_at = NOW() WHERE id = $1`,
			id, *req.Email,
		)
		if err != nil {
			jsonResponse(w, http.StatusInternalServerError, map[string]interface{}{
				"success": false, "message": "internal server error",
			})
			return
		}
		updated = true
	}

	if req.Password != nil {
		hash, err := bcrypt.GenerateFromPassword([]byte(*req.Password), bcrypt.DefaultCost)
		if err != nil {
			jsonResponse(w, http.StatusInternalServerError, map[string]interface{}{
				"success": false, "message": "internal server error",
			})
			return
		}
		_, err = h.Pool.Exec(r.Context(),
			`UPDATE users SET password_hash = $2, updated_at = NOW() WHERE id = $1`,
			id, string(hash),
		)
		if err != nil {
			jsonResponse(w, http.StatusInternalServerError, map[string]interface{}{
				"success": false, "message": "internal server error",
			})
			return
		}
		updated = true
	}

	if req.Credits != nil {
		_, err := h.Pool.Exec(r.Context(),
			`UPDATE users SET credits = $2, updated_at = NOW() WHERE id = $1`,
			id, *req.Credits,
		)
		if err != nil {
			jsonResponse(w, http.StatusInternalServerError, map[string]interface{}{
				"success": false, "message": "internal server error",
			})
			return
		}
		updated = true
	}

	if !updated {
		jsonResponse(w, http.StatusBadRequest, map[string]interface{}{
			"success": false, "message": "no fields to update",
		})
		return
	}

	// Create audit log for user update
	if h.AuditRepo != nil {
		_ = h.AuditRepo.Create(r.Context(), &audit.Log{
			UserID:    adminID,
			Action:    "update",
			Resource:  "user",
			Detail:    map[string]interface{}{"user_id": id, "fields": req},
			IPAddress: r.RemoteAddr,
		})
	}

	var usr struct {
		ID       string `json:"id"`
		Email    string `json:"email"`
		Username string `json:"username"`
		Role     string `json:"role"`
		Status   string `json:"status"`
		Credits  int    `json:"credits"`
	}
	err := h.Pool.QueryRow(r.Context(),
		`SELECT id, COALESCE(email,''), COALESCE(username,''), COALESCE(role,'user'), COALESCE(status,'active'), COALESCE(credits,0) FROM users WHERE id = $1`, id,
	).Scan(&usr.ID, &usr.Email, &usr.Username, &usr.Role, &usr.Status, &usr.Credits)
	if err != nil {
		jsonResponse(w, http.StatusOK, map[string]interface{}{
			"success": true,
			"message": "用户已更新",
		})
		return
	}

	jsonResponse(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"data":    usr,
	})
}
