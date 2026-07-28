package handler

import (
	"encoding/json"
	"net/http"

	"github.com/aige/api/internal/middleware"
	"github.com/aige/plan"
	"github.com/go-chi/chi/v5"
)

type AdminPlanHandler struct {
	PlanRepo plan.Repository
}

func (h *AdminPlanHandler) ListPlans(w http.ResponseWriter, r *http.Request) {
	adminID := middleware.GetAdminID(r.Context())
	if adminID == "" {
		jsonResponse(w, http.StatusUnauthorized, map[string]interface{}{
			"success": false, "message": "unauthorized",
		})
		return
	}

	plans, err := h.PlanRepo.List(r.Context())
	if err != nil {
		jsonResponse(w, http.StatusInternalServerError, map[string]interface{}{
			"success": false, "message": "internal server error",
		})
		return
	}

	jsonResponse(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"data":    plans,
	})
}

func (h *AdminPlanHandler) CreatePlan(w http.ResponseWriter, r *http.Request) {
	adminID := middleware.GetAdminID(r.Context())
	if adminID == "" {
		jsonResponse(w, http.StatusUnauthorized, map[string]interface{}{
			"success": false, "message": "unauthorized",
		})
		return
	}

	var req plan.CreatePlanRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		jsonResponse(w, http.StatusBadRequest, map[string]interface{}{
			"success": false, "message": "invalid request body",
		})
		return
	}

	if req.Name == "" || req.Code == "" {
		jsonResponse(w, http.StatusBadRequest, map[string]interface{}{
			"success": false, "message": "name and code are required",
		})
		return
	}

	p, err := h.PlanRepo.Create(r.Context(), req)
	if err != nil {
		jsonResponse(w, http.StatusInternalServerError, map[string]interface{}{
			"success": false, "message": err.Error(),
		})
		return
	}

	jsonResponse(w, http.StatusCreated, map[string]interface{}{
		"success": true,
		"data":    p,
	})
}

func (h *AdminPlanHandler) UpdatePlan(w http.ResponseWriter, r *http.Request) {
	adminID := middleware.GetAdminID(r.Context())
	if adminID == "" {
		jsonResponse(w, http.StatusUnauthorized, map[string]interface{}{
			"success": false, "message": "unauthorized",
		})
		return
	}

	id := chi.URLParam(r, "id")

	var req plan.UpdatePlanRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		jsonResponse(w, http.StatusBadRequest, map[string]interface{}{
			"success": false, "message": "invalid request body",
		})
		return
	}

	p, err := h.PlanRepo.Update(r.Context(), id, req)
	if err != nil {
		jsonResponse(w, http.StatusInternalServerError, map[string]interface{}{
			"success": false, "message": err.Error(),
		})
		return
	}
	if p == nil {
		jsonResponse(w, http.StatusNotFound, map[string]interface{}{
			"success": false, "message": "plan not found",
		})
		return
	}

	jsonResponse(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"data":    p,
	})
}
