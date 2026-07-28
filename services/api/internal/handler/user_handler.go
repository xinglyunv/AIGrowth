package handler

import (
	"encoding/json"
	"net/http"

	"github.com/aige/api/internal/middleware"
	"github.com/aige/auth"
	"github.com/aige/user"
)

// UserHandler handles user profile endpoints.
type UserHandler struct {
	UserRepo  user.Repository
	JWTSecret string
}

// GetMe handles GET /api/v1/users/me.
func (h *UserHandler) GetMe(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())
	if userID == "" {
		jsonResponse(w, http.StatusUnauthorized, map[string]interface{}{
			"success": false, "message": "unauthorized",
		})
		return
	}

	u, err := h.UserRepo.FindByID(r.Context(), userID)
	if err != nil {
		jsonResponse(w, http.StatusInternalServerError, map[string]interface{}{
			"success": false, "message": "internal server error",
		})
		return
	}
	if u == nil {
		jsonResponse(w, http.StatusNotFound, map[string]interface{}{
			"success": false, "message": "user not found",
		})
		return
	}

	u.PasswordHash = ""
	jsonResponse(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"data":    u,
	})
}

// UpdateMe handles PUT /api/v1/users/me.
func (h *UserHandler) UpdateMe(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())
	if userID == "" {
		jsonResponse(w, http.StatusUnauthorized, map[string]interface{}{
			"success": false, "message": "unauthorized",
		})
		return
	}

	var req user.UpdateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		jsonResponse(w, http.StatusBadRequest, map[string]interface{}{
			"success": false, "message": "invalid request body",
		})
		return
	}

	u, err := h.UserRepo.Update(r.Context(), userID, req)
	if err != nil {
		jsonResponse(w, http.StatusInternalServerError, map[string]interface{}{
			"success": false, "message": "internal server error",
		})
		return
	}

	u.PasswordHash = ""
	jsonResponse(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"data":    u,
	})
}

// ChangePassword handles PUT /api/v1/users/me/password.
func (h *UserHandler) ChangePassword(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())
	if userID == "" {
		jsonResponse(w, http.StatusUnauthorized, map[string]interface{}{
			"success": false, "message": "unauthorized",
		})
		return
	}

	var req struct {
		OldPassword string `json:"old_password"`
		NewPassword string `json:"new_password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		jsonResponse(w, http.StatusBadRequest, map[string]interface{}{
			"success": false, "message": "invalid request body",
		})
		return
	}

	if req.OldPassword == "" || req.NewPassword == "" {
		jsonResponse(w, http.StatusBadRequest, map[string]interface{}{
			"success": false, "message": "old_password and new_password are required",
		})
		return
	}
	if len(req.NewPassword) < 8 {
		jsonResponse(w, http.StatusBadRequest, map[string]interface{}{
			"success": false, "message": "new password must be at least 8 characters",
		})
		return
	}

	u, err := h.UserRepo.FindByID(r.Context(), userID)
	if err != nil {
		jsonResponse(w, http.StatusInternalServerError, map[string]interface{}{
			"success": false, "message": "internal server error",
		})
		return
	}
	if u == nil {
		jsonResponse(w, http.StatusNotFound, map[string]interface{}{
			"success": false, "message": "user not found",
		})
		return
	}

	if !auth.CheckPassword(u.PasswordHash, req.OldPassword) {
		jsonResponse(w, http.StatusBadRequest, map[string]interface{}{
			"success": false, "message": "current password is incorrect",
		})
		return
	}

	hash, err := auth.HashPassword(req.NewPassword)
	if err != nil {
		jsonResponse(w, http.StatusInternalServerError, map[string]interface{}{
			"success": false, "message": "internal server error",
		})
		return
	}

	if err := h.UserRepo.UpdatePassword(r.Context(), userID, hash); err != nil {
		jsonResponse(w, http.StatusInternalServerError, map[string]interface{}{
			"success": false, "message": "internal server error",
		})
		return
	}

	jsonResponse(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"message": "password changed successfully",
	})
}
