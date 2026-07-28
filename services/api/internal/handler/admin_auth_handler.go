package handler

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/aige/admin"
	"github.com/aige/api/internal/middleware"
	"github.com/golang-jwt/jwt/v5"
)

type AdminAuthHandler struct {
	AdminRepo admin.Repository
	JWTSecret string
}

func (h *AdminAuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req admin.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		jsonResponse(w, http.StatusBadRequest, map[string]interface{}{
			"success": false, "message": "invalid request body",
		})
		return
	}

	if req.Email == "" || req.Password == "" {
		jsonResponse(w, http.StatusBadRequest, map[string]interface{}{
			"success": false, "message": "email and password are required",
		})
		return
	}

	a, err := h.AdminRepo.FindByEmail(r.Context(), req.Email)
	if err != nil {
		jsonResponse(w, http.StatusInternalServerError, map[string]interface{}{
			"success": false, "message": "internal server error",
		})
		return
	}
	if a == nil || a.Status != "active" {
		jsonResponse(w, http.StatusUnauthorized, map[string]interface{}{
			"success": false, "message": "invalid credentials",
		})
		return
	}

	if !h.AdminRepo.VerifyPassword(req.Password, a.PasswordHash) {
		jsonResponse(w, http.StatusUnauthorized, map[string]interface{}{
			"success": false, "message": "invalid credentials",
		})
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"admin_id": a.ID,
		"role":     a.Role,
		"exp":      time.Now().Add(24 * time.Hour).Unix(),
		"iat":      time.Now().Unix(),
	})

	tokenStr, err := token.SignedString([]byte(h.JWTSecret))
	if err != nil {
		jsonResponse(w, http.StatusInternalServerError, map[string]interface{}{
			"success": false, "message": "internal server error",
		})
		return
	}

	a.PasswordHash = ""

	jsonResponse(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"data": map[string]interface{}{
			"token": tokenStr,
			"admin": a,
		},
	})
}

func (h *AdminAuthHandler) Me(w http.ResponseWriter, r *http.Request) {
	adminID := middleware.GetAdminID(r.Context())
	if adminID == "" {
		jsonResponse(w, http.StatusUnauthorized, map[string]interface{}{
			"success": false, "message": "unauthorized",
		})
		return
	}

	a, err := h.AdminRepo.FindByID(r.Context(), adminID)
	if err != nil || a == nil {
		jsonResponse(w, http.StatusInternalServerError, map[string]interface{}{
			"success": false, "message": "internal server error",
		})
		return
	}

	a.PasswordHash = ""

	jsonResponse(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"data":    a,
	})
}
