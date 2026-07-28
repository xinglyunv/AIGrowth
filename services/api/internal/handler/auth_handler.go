package handler

import (
	"encoding/json"
	"net/http"
	"regexp"
	"strconv"
	"strings"

	"github.com/aige/auth"
	"github.com/aige/setting"
	"github.com/aige/user"
)

// AuthHandler handles authentication-related endpoints.
type AuthHandler struct {
	UserRepo    user.Repository
	SettingRepo setting.Repository
	JWTSecret   string
}

// jsonResponse writes a JSON response.
func jsonResponse(w http.ResponseWriter, status int, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}

// Register handles POST /api/v1/auth/register.
func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Email       string  `json:"email"`
		Password    string  `json:"password"`
		Username    string  `json:"username"`
		CompanyName *string `json:"company_name,omitempty"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		jsonResponse(w, http.StatusBadRequest, map[string]interface{}{
			"success": false, "message": "invalid request body",
		})
		return
	}

	req.Email = strings.TrimSpace(req.Email)
	req.Username = strings.TrimSpace(req.Username)

	// Check registration toggle
	regSetting, err := h.SettingRepo.GetByKey(r.Context(), "allow_registration")
	if err == nil && regSetting != "" {
		enabled, parseErr := strconv.ParseBool(regSetting)
		if parseErr == nil && !enabled {
			jsonResponse(w, http.StatusForbidden, map[string]interface{}{
				"success": false, "message": "registration is currently disabled",
			})
			return
		}
	}

	// Validate
	if req.Email == "" || !isValidEmail(req.Email) {
		jsonResponse(w, http.StatusBadRequest, map[string]interface{}{
			"success": false, "message": "valid email is required",
		})
		return
	}
	if len(req.Password) < 8 {
		jsonResponse(w, http.StatusBadRequest, map[string]interface{}{
			"success": false, "message": "password must be at least 8 characters",
		})
		return
	}
	if req.Username == "" {
		jsonResponse(w, http.StatusBadRequest, map[string]interface{}{
			"success": false, "message": "username is required",
		})
		return
	}

	// Check email uniqueness
	existing, err := h.UserRepo.FindByEmail(r.Context(), req.Email)
	if err != nil {
		jsonResponse(w, http.StatusInternalServerError, map[string]interface{}{
			"success": false, "message": "internal server error",
		})
		return
	}
	if existing != nil {
		jsonResponse(w, http.StatusConflict, map[string]interface{}{
			"success": false, "message": "email already registered",
		})
		return
	}

	// Hash password
	hash, err := auth.HashPassword(req.Password)
	if err != nil {
		jsonResponse(w, http.StatusInternalServerError, map[string]interface{}{
			"success": false, "message": "internal server error",
		})
		return
	}

	// Create user
	u, err := h.UserRepo.Create(r.Context(), user.CreateUserRequest{
		Email:       req.Email,
		Password:    hash,
		Username:    req.Username,
		CompanyName: req.CompanyName,
	})
	if err != nil {
		jsonResponse(w, http.StatusInternalServerError, map[string]interface{}{
			"success": false, "message": "internal server error",
		})
		return
	}

	// Generate token
	token, err := auth.GenerateToken(u.ID, u.Role, h.JWTSecret)
	if err != nil {
		jsonResponse(w, http.StatusInternalServerError, map[string]interface{}{
			"success": false, "message": "internal server error",
		})
		return
	}

	jsonResponse(w, http.StatusCreated, map[string]interface{}{
		"success": true,
		"data": map[string]interface{}{
			"token": token,
			"user":  u,
		},
	})
}

// Login handles POST /api/v1/auth/login.
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Email    string `json:"email"`
		Phone    string `json:"phone"`
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		jsonResponse(w, http.StatusBadRequest, map[string]interface{}{
			"success": false, "message": "invalid request body",
		})
		return
	}

	req.Email = strings.TrimSpace(req.Email)
	req.Phone = strings.TrimSpace(req.Phone)
	if (req.Email == "" && req.Phone == "") || req.Password == "" {
		jsonResponse(w, http.StatusBadRequest, map[string]interface{}{
			"success": false, "message": "email or phone and password are required",
		})
		return
	}

	var u *user.User
	var findErr error

	if req.Phone != "" {
		u, findErr = h.UserRepo.FindByPhone(r.Context(), req.Phone)
	} else {
		u, findErr = h.UserRepo.FindByEmail(r.Context(), req.Email)
	}

	if findErr != nil {
		jsonResponse(w, http.StatusInternalServerError, map[string]interface{}{
			"success": false, "message": "internal server error",
		})
		return
	}
	if u == nil || !auth.CheckPassword(u.PasswordHash, req.Password) {
		jsonResponse(w, http.StatusUnauthorized, map[string]interface{}{
			"success": false, "message": "invalid credentials",
		})
		return
	}

	token, err := auth.GenerateToken(u.ID, u.Role, h.JWTSecret)
	if err != nil {
		jsonResponse(w, http.StatusInternalServerError, map[string]interface{}{
			"success": false, "message": "internal server error",
		})
		return
	}

	// Update last login (non-critical, ignore error)
	_ = h.UserRepo.UpdateLastLogin(r.Context(), u.ID)

	// Clear sensitive field for response
	u.PasswordHash = ""

	jsonResponse(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"data": map[string]interface{}{
			"token": token,
			"user":  u,
		},
	})
}

// SendCode handles POST /api/v1/auth/send-code.
func (h *AuthHandler) SendCode(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Email   string `json:"email"`
		Purpose string `json:"purpose"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		jsonResponse(w, http.StatusBadRequest, map[string]interface{}{
			"success": false, "message": "invalid request body",
		})
		return
	}

	req.Email = strings.TrimSpace(req.Email)
	if req.Email == "" || !isValidEmail(req.Email) {
		jsonResponse(w, http.StatusBadRequest, map[string]interface{}{
			"success": false, "message": "valid email is required",
		})
		return
	}
	if req.Purpose == "" {
		req.Purpose = "reset_password"
	}

	code := auth.GenerateVerificationCode()

	if err := h.UserRepo.SaveVerificationCode(r.Context(), req.Email, code, req.Purpose); err != nil {
		jsonResponse(w, http.StatusInternalServerError, map[string]interface{}{
			"success": false, "message": "internal server error",
		})
		return
	}

	// In production, send email here. For now, return the code in dev mode.
	resp := map[string]interface{}{
		"success": true,
		"message": "verification code sent",
	}
	jsonResponse(w, http.StatusOK, resp)
}

// ResetPassword handles POST /api/v1/auth/reset-password.
func (h *AuthHandler) ResetPassword(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Email       string `json:"email"`
		Code        string `json:"code"`
		NewPassword string `json:"new_password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		jsonResponse(w, http.StatusBadRequest, map[string]interface{}{
			"success": false, "message": "invalid request body",
		})
		return
	}

	req.Email = strings.TrimSpace(req.Email)
	if req.Email == "" || req.Code == "" || req.NewPassword == "" {
		jsonResponse(w, http.StatusBadRequest, map[string]interface{}{
			"success": false, "message": "email, code and new_password are required",
		})
		return
	}
	if len(req.NewPassword) < 8 {
		jsonResponse(w, http.StatusBadRequest, map[string]interface{}{
			"success": false, "message": "new password must be at least 8 characters",
		})
		return
	}

	// Verify code
	valid, err := h.UserRepo.VerifyCode(r.Context(), req.Email, req.Code, "reset_password")
	if err != nil {
		jsonResponse(w, http.StatusInternalServerError, map[string]interface{}{
			"success": false, "message": "internal server error",
		})
		return
	}
	if !valid {
		jsonResponse(w, http.StatusBadRequest, map[string]interface{}{
			"success": false, "message": "invalid or expired verification code",
		})
		return
	}

	// Find user
	u, err := h.UserRepo.FindByEmail(r.Context(), req.Email)
	if err != nil {
		jsonResponse(w, http.StatusInternalServerError, map[string]interface{}{
			"success": false, "message": "internal server error",
		})
		return
	}
	if u == nil {
		jsonResponse(w, http.StatusBadRequest, map[string]interface{}{
			"success": false, "message": "user not found",
		})
		return
	}

	// Hash and update password
	hash, err := auth.HashPassword(req.NewPassword)
	if err != nil {
		jsonResponse(w, http.StatusInternalServerError, map[string]interface{}{
			"success": false, "message": "internal server error",
		})
		return
	}

	if err := h.UserRepo.UpdatePassword(r.Context(), u.ID, hash); err != nil {
		jsonResponse(w, http.StatusInternalServerError, map[string]interface{}{
			"success": false, "message": "internal server error",
		})
		return
	}

	// Mark code as used (non-critical)
	_ = h.UserRepo.MarkCodeUsed(r.Context(), req.Email, req.Code, "reset_password")

	jsonResponse(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"message": "password reset successfully",
	})
}

var emailRegex = regexp.MustCompile(`^.+@.+\..+$`)

func isValidEmail(email string) bool {
	return emailRegex.MatchString(email)
}
