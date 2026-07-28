package handler

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/aige/api/internal/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type AdminPromptHandler struct {
	Pool *pgxpool.Pool
}

type Prompt struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Purpose   string    `json:"purpose"`
	Version   string    `json:"version"`
	Content   string    `json:"content"`
	Status    string    `json:"status"`
	CreatedBy string    `json:"created_by"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (h *AdminPromptHandler) ListPrompts(w http.ResponseWriter, r *http.Request) {
	adminID := middleware.GetAdminID(r.Context())
	if adminID == "" {
		jsonResponse(w, http.StatusUnauthorized, map[string]interface{}{
			"success": false, "message": "unauthorized",
		})
		return
	}

	rows, err := h.Pool.Query(r.Context(),
		`SELECT id, name, purpose, COALESCE(version, 'v1.0') AS version, content, status,
		        COALESCE(created_by::text, '') AS created_by, created_at, updated_at
		 FROM prompts ORDER BY created_at DESC`,
	)
	if err != nil {
		jsonResponse(w, http.StatusInternalServerError, map[string]interface{}{
			"success": false, "message": "internal server error",
		})
		return
	}
	defer rows.Close()

	var prompts []Prompt
	for rows.Next() {
		var p Prompt
		err := rows.Scan(&p.ID, &p.Name, &p.Purpose, &p.Version, &p.Content,
			&p.Status, &p.CreatedBy, &p.CreatedAt, &p.UpdatedAt)
		if err != nil {
			jsonResponse(w, http.StatusInternalServerError, map[string]interface{}{
				"success": false, "message": "internal server error",
			})
			return
		}
		prompts = append(prompts, p)
	}
	if prompts == nil {
		prompts = []Prompt{}
	}

	jsonResponse(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"data":    prompts,
	})
}

func (h *AdminPromptHandler) CreatePrompt(w http.ResponseWriter, r *http.Request) {
	adminID := middleware.GetAdminID(r.Context())
	if adminID == "" {
		jsonResponse(w, http.StatusUnauthorized, map[string]interface{}{
			"success": false, "message": "unauthorized",
		})
		return
	}

	var req struct {
		Name    string `json:"name"`
		Purpose string `json:"purpose"`
		Version string `json:"version"`
		Content string `json:"content"`
		Status  string `json:"status"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		jsonResponse(w, http.StatusBadRequest, map[string]interface{}{
			"success": false, "message": "invalid request body",
		})
		return
	}

	if req.Name == "" || req.Purpose == "" || req.Content == "" {
		jsonResponse(w, http.StatusBadRequest, map[string]interface{}{
			"success": false, "message": "name, purpose, and content are required",
		})
		return
	}

	if req.Version == "" {
		req.Version = "v1.0"
	}
	if req.Status == "" {
		req.Status = "draft"
	}

	var p Prompt
	err := h.Pool.QueryRow(r.Context(),
		`INSERT INTO prompts (name, purpose, version, content, status, created_by)
		 VALUES ($1, $2, $3, $4, $5, $6)
		 RETURNING id, name, purpose, version, content, status, COALESCE(created_by::text, '') AS created_by, created_at, updated_at`,
		req.Name, req.Purpose, req.Version, req.Content, req.Status, adminID,
	).Scan(&p.ID, &p.Name, &p.Purpose, &p.Version, &p.Content, &p.Status, &p.CreatedBy, &p.CreatedAt, &p.UpdatedAt)
	if err != nil {
		jsonResponse(w, http.StatusInternalServerError, map[string]interface{}{
			"success": false, "message": "internal server error",
		})
		return
	}

	jsonResponse(w, http.StatusCreated, map[string]interface{}{
		"success": true,
		"data":    p,
	})
}

func (h *AdminPromptHandler) GetPrompt(w http.ResponseWriter, r *http.Request) {
	adminID := middleware.GetAdminID(r.Context())
	if adminID == "" {
		jsonResponse(w, http.StatusUnauthorized, map[string]interface{}{
			"success": false, "message": "unauthorized",
		})
		return
	}

	id := chi.URLParam(r, "id")

	var p Prompt
	err := h.Pool.QueryRow(r.Context(),
		`SELECT id, name, purpose, COALESCE(version, 'v1.0') AS version, content, status,
		        COALESCE(created_by::text, '') AS created_by, created_at, updated_at
		 FROM prompts WHERE id = $1`, id,
	).Scan(&p.ID, &p.Name, &p.Purpose, &p.Version, &p.Content, &p.Status, &p.CreatedBy, &p.CreatedAt, &p.UpdatedAt)
	if err != nil {
		jsonResponse(w, http.StatusNotFound, map[string]interface{}{
			"success": false, "message": "prompt not found",
		})
		return
	}

	jsonResponse(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"data":    p,
	})
}

func (h *AdminPromptHandler) UpdatePrompt(w http.ResponseWriter, r *http.Request) {
	adminID := middleware.GetAdminID(r.Context())
	if adminID == "" {
		jsonResponse(w, http.StatusUnauthorized, map[string]interface{}{
			"success": false, "message": "unauthorized",
		})
		return
	}

	id := chi.URLParam(r, "id")

	var req struct {
		Name    *string `json:"name"`
		Purpose *string `json:"purpose"`
		Version *string `json:"version"`
		Content *string `json:"content"`
		Status  *string `json:"status"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		jsonResponse(w, http.StatusBadRequest, map[string]interface{}{
			"success": false, "message": "invalid request body",
		})
		return
	}

	setClauses := []string{}
	args := []interface{}{}
	argIdx := 1

	if req.Name != nil {
		setClauses = append(setClauses, "name = $"+strconv.Itoa(argIdx))
		args = append(args, *req.Name)
		argIdx++
	}
	if req.Purpose != nil {
		setClauses = append(setClauses, "purpose = $"+strconv.Itoa(argIdx))
		args = append(args, *req.Purpose)
		argIdx++
	}
	if req.Version != nil {
		setClauses = append(setClauses, "version = $"+strconv.Itoa(argIdx))
		args = append(args, *req.Version)
		argIdx++
	}
	if req.Content != nil {
		setClauses = append(setClauses, "content = $"+strconv.Itoa(argIdx))
		args = append(args, *req.Content)
		argIdx++
	}
	if req.Status != nil {
		setClauses = append(setClauses, "status = $"+strconv.Itoa(argIdx))
		args = append(args, *req.Status)
		argIdx++
	}

	if len(setClauses) == 0 {
		jsonResponse(w, http.StatusBadRequest, map[string]interface{}{
			"success": false, "message": "no fields to update",
		})
		return
	}

	setClauses = append(setClauses, "updated_at = NOW()")
	args = append(args, id)

	query := `UPDATE prompts SET ` + joinStrings(setClauses, ", ") + ` WHERE id = $` + strconv.Itoa(argIdx) +
		` RETURNING id, name, purpose, COALESCE(version, 'v1.0') AS version, content, status,
		           COALESCE(created_by::text, '') AS created_by, created_at, updated_at`

	var p Prompt
	err := h.Pool.QueryRow(r.Context(), query, args...).Scan(
		&p.ID, &p.Name, &p.Purpose, &p.Version, &p.Content, &p.Status, &p.CreatedBy, &p.CreatedAt, &p.UpdatedAt,
	)
	if err != nil {
		jsonResponse(w, http.StatusNotFound, map[string]interface{}{
			"success": false, "message": "prompt not found",
		})
		return
	}

	jsonResponse(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"data":    p,
	})
}

func (h *AdminPromptHandler) PublishPrompt(w http.ResponseWriter, r *http.Request) {
	adminID := middleware.GetAdminID(r.Context())
	if adminID == "" {
		jsonResponse(w, http.StatusUnauthorized, map[string]interface{}{
			"success": false, "message": "unauthorized",
		})
		return
	}

	id := chi.URLParam(r, "id")

	var p Prompt
	err := h.Pool.QueryRow(r.Context(),
		`UPDATE prompts SET status = 'published', updated_at = NOW()
		 WHERE id = $1
		 RETURNING id, name, purpose, COALESCE(version, 'v1.0') AS version, content, status,
		           COALESCE(created_by::text, '') AS created_by, created_at, updated_at`,
		id,
	).Scan(&p.ID, &p.Name, &p.Purpose, &p.Version, &p.Content, &p.Status, &p.CreatedBy, &p.CreatedAt, &p.UpdatedAt)
	if err != nil {
		jsonResponse(w, http.StatusNotFound, map[string]interface{}{
			"success": false, "message": "prompt not found",
		})
		return
	}

	jsonResponse(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"data":    p,
		"message": "提示词已发布",
	})
}

func joinStrings(elems []string, sep string) string {
	if len(elems) == 0 {
		return ""
	}
	result := elems[0]
	for _, e := range elems[1:] {
		result += sep + e
	}
	return result
}
