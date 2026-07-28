package handler

import (
	"log"
	"net/http"
	"strconv"

	"github.com/aige/api/internal/middleware"
	"github.com/aige/report"
	"github.com/go-chi/chi/v5"
)

type ReportHandler struct {
	ReportRepo report.Repository
}

func (h *ReportHandler) List(w http.ResponseWriter, r *http.Request) {
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

	reports, total, err := h.ReportRepo.ListByUserID(r.Context(), userID, offset, limit)
	if err != nil {
		log.Printf("ERROR listing reports: %v", err)
		jsonResponse(w, http.StatusInternalServerError, map[string]interface{}{
			"success": false, "message": "internal server error",
		})
		return
	}

	jsonResponse(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"data":    reports,
		"meta": map[string]interface{}{
			"total":  total,
			"offset": offset,
			"limit":  limit,
		},
	})
}

func (h *ReportHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())
	if userID == "" {
		jsonResponse(w, http.StatusUnauthorized, map[string]interface{}{
			"success": false, "message": "unauthorized",
		})
		return
	}

	id := chi.URLParam(r, "id")
	report, err := h.ReportRepo.GetByID(r.Context(), id)
	if err != nil {
		log.Printf("ERROR getting report: %v", err)
		jsonResponse(w, http.StatusInternalServerError, map[string]interface{}{
			"success": false, "message": "internal server error",
		})
		return
	}
	if report == nil {
		jsonResponse(w, http.StatusNotFound, map[string]interface{}{
			"success": false, "message": "report not found",
		})
		return
	}
	if report.UserID != userID {
		jsonResponse(w, http.StatusForbidden, map[string]interface{}{
			"success": false, "message": "access denied",
		})
		return
	}

	jsonResponse(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"data":    report,
	})
}

func (h *ReportHandler) ListByProject(w http.ResponseWriter, r *http.Request) {
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

	offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	if limit <= 0 || limit > 100 {
		limit = 20
	}
	if offset < 0 {
		offset = 0
	}

	reports, total, err := h.ReportRepo.ListByProjectID(r.Context(), projectID, offset, limit)
	if err != nil {
		log.Printf("ERROR listing reports by project: %v", err)
		jsonResponse(w, http.StatusInternalServerError, map[string]interface{}{
			"success": false, "message": "internal server error",
		})
		return
	}

	jsonResponse(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"data":    reports,
		"meta": map[string]interface{}{
			"total":  total,
			"offset": offset,
			"limit":  limit,
		},
	})
}
