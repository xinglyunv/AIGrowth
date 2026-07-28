package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/aige/admin"
	"github.com/go-chi/chi/v5"
)

type AdminHandler struct {
	AdminRepo admin.Repository
}

func (h *AdminHandler) List(w http.ResponseWriter, r *http.Request) {
	offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	if limit <= 0 || limit > 100 {
		limit = 20
	}
	if offset < 0 {
		offset = 0
	}

	admins, total, err := h.AdminRepo.List(r.Context(), offset, limit)
	if err != nil {
		jsonResponse(w, http.StatusInternalServerError, map[string]interface{}{
			"success": false, "message": "internal server error",
		})
		return
	}

	for i := range admins {
		admins[i].PasswordHash = ""
	}

	jsonResponse(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"data":    admins,
		"meta": map[string]interface{}{
			"total":  total,
			"offset": offset,
			"limit":  limit,
		},
	})
}

func (h *AdminHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req admin.CreateAdminRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		jsonResponse(w, http.StatusBadRequest, map[string]interface{}{
			"success": false, "message": "invalid request body",
		})
		return
	}

	if req.Username == "" || req.Email == "" || req.Password == "" {
		jsonResponse(w, http.StatusBadRequest, map[string]interface{}{
			"success": false, "message": "username, email, and password are required",
		})
		return
	}

	a, err := h.AdminRepo.Create(r.Context(), req)
	if err != nil {
		jsonResponse(w, http.StatusInternalServerError, map[string]interface{}{
			"success": false, "message": err.Error(),
		})
		return
	}

	a.PasswordHash = ""

	jsonResponse(w, http.StatusCreated, map[string]interface{}{
		"success": true,
		"data":    a,
	})
}

func (h *AdminHandler) Get(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	a, err := h.AdminRepo.FindByID(r.Context(), id)
	if err != nil || a == nil {
		jsonResponse(w, http.StatusNotFound, map[string]interface{}{
			"success": false, "message": "admin not found",
		})
		return
	}

	a.PasswordHash = ""

	jsonResponse(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"data":    a,
	})
}

func (h *AdminHandler) Update(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var req admin.UpdateAdminRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		jsonResponse(w, http.StatusBadRequest, map[string]interface{}{
			"success": false, "message": "invalid request body",
		})
		return
	}

	a, err := h.AdminRepo.Update(r.Context(), id, req)
	if err != nil {
		jsonResponse(w, http.StatusInternalServerError, map[string]interface{}{
			"success": false, "message": err.Error(),
		})
		return
	}
	if a == nil {
		jsonResponse(w, http.StatusNotFound, map[string]interface{}{
			"success": false, "message": "admin not found",
		})
		return
	}

	a.PasswordHash = ""

	jsonResponse(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"data":    a,
	})
}

func (h *AdminHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	err := h.AdminRepo.Delete(r.Context(), id)
	if err != nil {
		jsonResponse(w, http.StatusInternalServerError, map[string]interface{}{
			"success": false, "message": "internal server error",
		})
		return
	}

	jsonResponse(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"message": "admin deleted",
	})
}
