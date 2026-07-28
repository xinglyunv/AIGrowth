package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/mail"
	"strings"

	"github.com/aige/api/internal/middleware"
	"github.com/aige/space"
	"github.com/go-chi/chi/v5"
)

type SpaceHandler struct{ SpaceRepo space.Repository }

func (h *SpaceHandler) List(w http.ResponseWriter, r *http.Request) {
	spaces, err := h.SpaceRepo.ListByUser(r.Context(), middleware.GetUserID(r.Context()))
	if err != nil {
		jsonResponse(w, http.StatusInternalServerError, map[string]interface{}{"success": false, "message": "internal server error"})
		return
	}
	jsonResponse(w, http.StatusOK, map[string]interface{}{"success": true, "data": spaces})
}

func (h *SpaceHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req space.CreateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		jsonResponse(w, http.StatusBadRequest, map[string]interface{}{"success": false, "message": "invalid request body"})
		return
	}
	req.Name = strings.TrimSpace(req.Name)
	if req.Name == "" || len([]rune(req.Name)) > 120 {
		jsonResponse(w, http.StatusBadRequest, map[string]interface{}{"success": false, "message": "workspace name must be between 1 and 120 characters"})
		return
	}
	s, err := h.SpaceRepo.Create(r.Context(), middleware.GetUserID(r.Context()), req)
	if err != nil {
		jsonResponse(w, http.StatusInternalServerError, map[string]interface{}{"success": false, "message": "internal server error"})
		return
	}
	jsonResponse(w, http.StatusCreated, map[string]interface{}{"success": true, "data": s})
}

func (h *SpaceHandler) Current(w http.ResponseWriter, r *http.Request) {
	s, err := h.SpaceRepo.GetCurrent(r.Context(), middleware.GetUserID(r.Context()))
	if err != nil {
		jsonResponse(w, http.StatusInternalServerError, map[string]interface{}{"success": false, "message": "internal server error"})
		return
	}
	if s == nil {
		jsonResponse(w, http.StatusNotFound, map[string]interface{}{"success": false, "message": "current workspace not found"})
		return
	}
	jsonResponse(w, http.StatusOK, map[string]interface{}{"success": true, "data": s})
}

func (h *SpaceHandler) SetCurrent(w http.ResponseWriter, r *http.Request) {
	s, err := h.SpaceRepo.SetCurrent(r.Context(), middleware.GetUserID(r.Context()), chi.URLParam(r, "id"))
	if err != nil {
		jsonResponse(w, http.StatusInternalServerError, map[string]interface{}{"success": false, "message": "internal server error"})
		return
	}
	if s == nil {
		jsonResponse(w, http.StatusForbidden, map[string]interface{}{"success": false, "message": "workspace access denied"})
		return
	}
	jsonResponse(w, http.StatusOK, map[string]interface{}{"success": true, "data": s})
}

func (h *SpaceHandler) Members(w http.ResponseWriter, r *http.Request) {
	members, err := h.SpaceRepo.ListMembers(r.Context(), middleware.GetUserID(r.Context()), chi.URLParam(r, "id"))
	if err != nil {
		jsonResponse(w, http.StatusInternalServerError, map[string]interface{}{"success": false, "message": "internal server error"})
		return
	}
	jsonResponse(w, http.StatusOK, map[string]interface{}{"success": true, "data": members})
}

func (h *SpaceHandler) Invite(w http.ResponseWriter, r *http.Request) {
	var req space.InviteRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		jsonResponse(w, http.StatusBadRequest, map[string]interface{}{"success": false, "message": "invalid request body"})
		return
	}
	req.Email = strings.ToLower(strings.TrimSpace(req.Email))
	if parsed, err := mail.ParseAddress(req.Email); err != nil || parsed.Address != req.Email {
		jsonResponse(w, http.StatusBadRequest, map[string]interface{}{"success": false, "message": "valid email is required"})
		return
	}
	if req.Role == "" {
		req.Role = space.RoleMember
	}
	if !space.IsValidRole(req.Role, false) {
		jsonResponse(w, http.StatusBadRequest, map[string]interface{}{"success": false, "message": "invalid workspace role"})
		return
	}
	ok, err := h.SpaceRepo.Invite(r.Context(), middleware.GetUserID(r.Context()), chi.URLParam(r, "id"), req)
	if err != nil {
		jsonResponse(w, http.StatusInternalServerError, map[string]interface{}{"success": false, "message": "internal server error"})
		return
	}
	if !ok {
		jsonResponse(w, http.StatusForbidden, map[string]interface{}{"success": false, "message": "workspace management permission required"})
		return
	}
	inviteURL := fmt.Sprintf("https://%s/app/spaces/%s", r.Host, chi.URLParam(r, "id"))
	jsonResponse(w, http.StatusCreated, map[string]interface{}{"success": true, "message": "invitation created", "data": map[string]interface{}{"invite_url": inviteURL}})
}

func (h *SpaceHandler) UpdateMemberRole(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Role string `json:"role"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || !space.IsValidRole(req.Role, false) {
		jsonResponse(w, http.StatusBadRequest, map[string]interface{}{"success": false, "message": "invalid workspace role"})
		return
	}
	ok, err := h.SpaceRepo.UpdateMemberRole(r.Context(), middleware.GetUserID(r.Context()), chi.URLParam(r, "id"), chi.URLParam(r, "userID"), req.Role)
	if err != nil {
		jsonResponse(w, http.StatusInternalServerError, map[string]interface{}{"success": false, "message": "internal server error"})
		return
	}
	if !ok {
		jsonResponse(w, http.StatusForbidden, map[string]interface{}{"success": false, "message": "member role cannot be changed"})
		return
	}
	jsonResponse(w, http.StatusOK, map[string]interface{}{"success": true})
}

func (h *SpaceHandler) RemoveMember(w http.ResponseWriter, r *http.Request) {
	ok, err := h.SpaceRepo.RemoveMember(r.Context(), middleware.GetUserID(r.Context()), chi.URLParam(r, "id"), chi.URLParam(r, "userID"))
	if err != nil {
		jsonResponse(w, http.StatusInternalServerError, map[string]interface{}{"success": false, "message": "internal server error"})
		return
	}
	if !ok {
		jsonResponse(w, http.StatusForbidden, map[string]interface{}{"success": false, "message": "member cannot be removed"})
		return
	}
	jsonResponse(w, http.StatusOK, map[string]interface{}{"success": true})
}
