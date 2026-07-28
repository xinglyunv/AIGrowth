package handler

import (
	"log"
	"net/http"

	"github.com/aige/api/internal/middleware"
	"github.com/aige/competitor"
	"github.com/aige/project"
	"github.com/go-chi/chi/v5"
)

type CompetitorHandler struct {
	CompRepo    competitor.Repository
	ProjectRepo project.Repository
}

func (h *CompetitorHandler) ListByProject(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())
	if userID == "" {
		jsonResponse(w, http.StatusUnauthorized, map[string]interface{}{
			"success": false, "message": "unauthorized",
		})
		return
	}

	projectID := chi.URLParam(r, "projectId")
	if projectID == "" {
		jsonResponse(w, http.StatusBadRequest, map[string]interface{}{
			"success": false, "message": "project id is required",
		})
		return
	}
	projectRecord, err := h.ProjectRepo.FindByID(r.Context(), projectID)
	if err != nil {
		log.Printf("ERROR finding project for competitors: %v", err)
		jsonResponse(w, http.StatusInternalServerError, map[string]interface{}{
			"success": false, "message": "internal server error",
		})
		return
	}
	if projectRecord == nil || projectRecord.UserID != userID {
		jsonResponse(w, http.StatusForbidden, map[string]interface{}{
			"success": false, "message": "access denied",
		})
		return
	}

	competitors, err := h.CompRepo.ListByProjectID(r.Context(), projectID)
	if err != nil {
		log.Printf("ERROR listing competitors: %v", err)
		jsonResponse(w, http.StatusInternalServerError, map[string]interface{}{
			"success": false, "message": "internal server error",
		})
		return
	}

	jsonResponse(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"data":    competitors,
	})
}
