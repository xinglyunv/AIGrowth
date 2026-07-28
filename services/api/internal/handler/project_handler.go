package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/aige/aimodel"
	"github.com/aige/api/internal/middleware"
	"github.com/aige/project"
	"github.com/aige/user"
	"github.com/go-chi/chi/v5"
)

// ProjectHandler handles brand project endpoints.
type ProjectHandler struct {
	ProjectRepo project.Repository
	ModelRepo   aimodel.Repository
	AIProvider  aimodel.AIProvider
	UserRepo    user.Repository
}

// Create handles POST /api/v1/projects.
func (h *ProjectHandler) Create(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())
	if userID == "" {
		jsonResponse(w, http.StatusUnauthorized, map[string]interface{}{
			"success": false, "message": "unauthorized",
		})
		return
	}

	var req project.CreateProjectRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		jsonResponse(w, http.StatusBadRequest, map[string]interface{}{
			"success": false, "message": "invalid request body",
		})
		return
	}

	req.Name = strings.TrimSpace(req.Name)
	req.Industry = strings.TrimSpace(req.Industry)

	if req.Name == "" {
		jsonResponse(w, http.StatusBadRequest, map[string]interface{}{
			"success": false, "message": "project name is required",
		})
		return
	}
	if req.Industry == "" {
		jsonResponse(w, http.StatusBadRequest, map[string]interface{}{
			"success": false, "message": "industry is required",
		})
		return
	}

	p, err := h.ProjectRepo.Create(r.Context(), userID, req)
	if err != nil {
		log.Printf("ERROR creating project: %v", err)
		jsonResponse(w, http.StatusInternalServerError, map[string]interface{}{
			"success": false, "message": "internal server error",
		})
		return
	}

	// Background brand info generation via AI
	go h.generateBrandInfo(context.Background(), p.ID, p.Name, p.Industry, p.Keywords)

	jsonResponse(w, http.StatusCreated, map[string]interface{}{
		"success": true,
		"data":    p,
	})
}

// List handles GET /api/v1/projects.
func (h *ProjectHandler) List(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())
	if userID == "" {
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

	projects, total, err := h.ProjectRepo.ListByUser(r.Context(), userID, offset, limit)
	if err != nil {
		log.Printf("ERROR listing projects: %v", err)
		jsonResponse(w, http.StatusInternalServerError, map[string]interface{}{
			"success": false, "message": "internal server error",
		})
		return
	}

	jsonResponse(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"data":    projects,
		"meta": map[string]interface{}{
			"total":  total,
			"offset": offset,
			"limit":  limit,
		},
	})
}

// Get handles GET /api/v1/projects/{id}.
func (h *ProjectHandler) Get(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())
	if userID == "" {
		jsonResponse(w, http.StatusUnauthorized, map[string]interface{}{
			"success": false, "message": "unauthorized",
		})
		return
	}

	id := chi.URLParam(r, "id")
	p, err := h.ProjectRepo.FindByID(r.Context(), id)
	if err != nil {
		jsonResponse(w, http.StatusInternalServerError, map[string]interface{}{
			"success": false, "message": "internal server error",
		})
		return
	}
	if p == nil {
		jsonResponse(w, http.StatusNotFound, map[string]interface{}{
			"success": false, "message": "project not found",
		})
		return
	}
	if p.UserID != userID {
		jsonResponse(w, http.StatusForbidden, map[string]interface{}{
			"success": false, "message": "access denied",
		})
		return
	}

	jsonResponse(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"data":    p,
	})
}

// Update handles PUT /api/v1/projects/{id}.
func (h *ProjectHandler) Update(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())
	if userID == "" {
		jsonResponse(w, http.StatusUnauthorized, map[string]interface{}{
			"success": false, "message": "unauthorized",
		})
		return
	}

	id := chi.URLParam(r, "id")

	// Verify ownership
	existing, err := h.ProjectRepo.FindByID(r.Context(), id)
	if err != nil {
		jsonResponse(w, http.StatusInternalServerError, map[string]interface{}{
			"success": false, "message": "internal server error",
		})
		return
	}
	if existing == nil {
		jsonResponse(w, http.StatusNotFound, map[string]interface{}{
			"success": false, "message": "project not found",
		})
		return
	}
	if existing.UserID != userID {
		jsonResponse(w, http.StatusForbidden, map[string]interface{}{
			"success": false, "message": "access denied",
		})
		return
	}

	var req project.UpdateProjectRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		jsonResponse(w, http.StatusBadRequest, map[string]interface{}{
			"success": false, "message": "invalid request body",
		})
		return
	}

	p, err := h.ProjectRepo.Update(r.Context(), id, req)
	if err != nil {
		jsonResponse(w, http.StatusInternalServerError, map[string]interface{}{
			"success": false, "message": "internal server error",
		})
		return
	}
	if p == nil {
		jsonResponse(w, http.StatusNotFound, map[string]interface{}{
			"success": false, "message": "project not found",
		})
		return
	}

	// Upsert brand_info if any brand fields are provided
	if req.BrandIntro != nil || req.ProductIntro != nil || req.ServiceIntro != nil ||
		req.FAQ != nil || req.Advantages != nil || req.Cases != nil {
		if err := h.ProjectRepo.UpsertBrandInfo(r.Context(), id,
			req.BrandIntro, req.ProductIntro, req.ServiceIntro,
			req.FAQ, req.Advantages, req.Cases,
		); err != nil {
			log.Printf("ERROR upserting brand info: %v", err)
		} else {
			// Re-fetch to include updated brand_info
			p, _ = h.ProjectRepo.FindByID(r.Context(), id)
		}
	}

	jsonResponse(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"data":    p,
	})
}

// Delete handles DELETE /api/v1/projects/{id}.
func (h *ProjectHandler) Delete(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())
	if userID == "" {
		jsonResponse(w, http.StatusUnauthorized, map[string]interface{}{
			"success": false, "message": "unauthorized",
		})
		return
	}

	id := chi.URLParam(r, "id")

	// Verify ownership
	existing, err := h.ProjectRepo.FindByID(r.Context(), id)
	if err != nil {
		jsonResponse(w, http.StatusInternalServerError, map[string]interface{}{
			"success": false, "message": "internal server error",
		})
		return
	}
	if existing == nil {
		jsonResponse(w, http.StatusNotFound, map[string]interface{}{
			"success": false, "message": "project not found",
		})
		return
	}
	if existing.UserID != userID {
		jsonResponse(w, http.StatusForbidden, map[string]interface{}{
			"success": false, "message": "access denied",
		})
		return
	}

	if err := h.ProjectRepo.Delete(r.Context(), id); err != nil {
		jsonResponse(w, http.StatusInternalServerError, map[string]interface{}{
			"success": false, "message": "internal server error",
		})
		return
	}

	jsonResponse(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"message": "project archived successfully",
	})
}

// generateBrandInfo calls the AI provider to generate brand content.
func (h *ProjectHandler) generateBrandInfo(ctx context.Context, projectID, name, industry string, keywords []string) {
	ctx, cancel := context.WithTimeout(ctx, 60*time.Second)
	defer cancel()

	models, err := h.ModelRepo.ListEnabled(ctx)
	if err != nil || len(models) == 0 {
		log.Printf("ERROR no enabled AI model for brand generation: %v", err)
		return
	}

	keywordStr := strings.Join(keywords, ", ")
	prompt := fmt.Sprintf(`You are a brand analysis expert. Based on the following brand information, generate a JSON object with brand introduction, product introduction, and service introduction in Chinese.

Brand name: %s
Industry: %s
Keywords: %s

Output ONLY valid JSON (no markdown formatting) with these exact keys:
{
  "brand_intro": "...",
  "product_intro": "...",
  "service_intro": "..."
}`, name, industry, keywordStr)

	messages := []aimodel.ChatMessage{
		{Role: "system", Content: "You are a helpful assistant that outputs valid JSON only."},
		{Role: "user", Content: prompt},
	}

	content, err := h.AIProvider.Chat(ctx, models[0], messages)
	if err != nil {
		log.Printf("ERROR generating brand info for project %s: %v", projectID, err)
		return
	}

	content = cleanJSONResponse(content)

	var result struct {
		BrandIntro   string `json:"brand_intro"`
		ProductIntro string `json:"product_intro"`
		ServiceIntro string `json:"service_intro"`
	}
	if err := json.Unmarshal([]byte(content), &result); err != nil {
		log.Printf("ERROR parsing brand info JSON for project %s: %v", projectID, err)
		return
	}

	if err := h.ProjectRepo.UpsertBrandInfo(ctx, projectID,
		&result.BrandIntro, &result.ProductIntro, &result.ServiceIntro,
		nil, nil, nil,
	); err != nil {
		log.Printf("ERROR upserting generated brand info for project %s: %v", projectID, err)
		return
	}
	log.Printf("Brand info generated for project %s", projectID)
}

// cleanJSONResponse strips markdown code block fences from AI JSON output.
func cleanJSONResponse(s string) string {
	s = strings.TrimSpace(s)
	if strings.HasPrefix(s, "```") {
		if idx := strings.Index(s, "\n"); idx != -1 {
			s = s[idx+1:]
		}
		if idx := strings.LastIndex(s, "```"); idx != -1 {
			s = s[:idx]
		}
		s = strings.TrimSpace(s)
	}
	return s
}
