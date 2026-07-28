package handler

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/aige/api/internal/middleware"
	"github.com/aige/cdk"
	"github.com/go-chi/chi/v5"
)

type AdminCDKHandler struct {
	CDKRepo cdk.Repository
}

func (h *AdminCDKHandler) ListCDK(w http.ResponseWriter, r *http.Request) {
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

	codes, total, err := h.CDKRepo.List(r.Context(), offset, limit)
	if err != nil {
		jsonResponse(w, http.StatusInternalServerError, map[string]interface{}{
			"success": false, "message": "internal server error",
		})
		return
	}

	jsonResponse(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"data": map[string]interface{}{
			"data":  codes,
			"total": total,
		},
	})
}

func (h *AdminCDKHandler) CreateCDK(w http.ResponseWriter, r *http.Request) {
	adminID := middleware.GetAdminID(r.Context())
	if adminID == "" {
		jsonResponse(w, http.StatusUnauthorized, map[string]interface{}{
			"success": false, "message": "unauthorized",
		})
		return
	}

	var req cdk.CreateCDKRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		jsonResponse(w, http.StatusBadRequest, map[string]interface{}{
			"success": false, "message": "invalid request body",
		})
		return
	}

	c, err := h.CDKRepo.Create(r.Context(), req)
	if err != nil {
		if strings.Contains(err.Error(), "duplicate key") || strings.Contains(err.Error(), "unique constraint") {
			jsonResponse(w, http.StatusConflict, map[string]interface{}{"success": false, "message": "CDK 代码已存在"})
			return
		}
		jsonResponse(w, http.StatusInternalServerError, map[string]interface{}{
			"success": false, "message": err.Error(),
		})
		return
	}

	jsonResponse(w, http.StatusCreated, map[string]interface{}{
		"success": true,
		"data":    c,
	})
}

func (h *AdminCDKHandler) UpdateCDK(w http.ResponseWriter, r *http.Request) {
	adminID := middleware.GetAdminID(r.Context())
	if adminID == "" {
		jsonResponse(w, http.StatusUnauthorized, map[string]interface{}{
			"success": false, "message": "unauthorized",
		})
		return
	}

	id := chi.URLParam(r, "id")

	var req cdk.UpdateCDKRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		jsonResponse(w, http.StatusBadRequest, map[string]interface{}{
			"success": false, "message": "invalid request body",
		})
		return
	}

	c, err := h.CDKRepo.Update(r.Context(), id, req)
	if err != nil {
		jsonResponse(w, http.StatusInternalServerError, map[string]interface{}{
			"success": false, "message": err.Error(),
		})
		return
	}
	if c == nil {
		jsonResponse(w, http.StatusNotFound, map[string]interface{}{
			"success": false, "message": "CDK not found",
		})
		return
	}

	jsonResponse(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"data":    c,
	})
}

func (h *AdminCDKHandler) GetCDKUsages(w http.ResponseWriter, r *http.Request) {
	adminID := middleware.GetAdminID(r.Context())
	if adminID == "" {
		jsonResponse(w, http.StatusUnauthorized, map[string]interface{}{
			"success": false, "message": "unauthorized",
		})
		return
	}

	id := chi.URLParam(r, "id")

	usages, err := h.CDKRepo.GetUsages(r.Context(), id)
	if err != nil {
		jsonResponse(w, http.StatusInternalServerError, map[string]interface{}{
			"success": false, "message": "internal server error",
		})
		return
	}

	jsonResponse(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"data":    usages,
	})
}
