package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/aige/contact"
	"github.com/go-chi/chi/v5"
)

type ContactHandler struct {
	ContactRepo contact.Repository
}

func (h *ContactHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req contact.CreateMessageRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		jsonResponse(w, http.StatusBadRequest, map[string]interface{}{
			"success": false, "message": "invalid request body",
		})
		return
	}
	if req.Name == "" || req.Email == "" || req.Message == "" {
		jsonResponse(w, http.StatusBadRequest, map[string]interface{}{
			"success": false, "message": "name, email and message are required",
		})
		return
	}
	if err := h.ContactRepo.Create(r.Context(), &req); err != nil {
		jsonResponse(w, http.StatusInternalServerError, map[string]interface{}{
			"success": false, "message": "failed to create message",
		})
		return
	}
	jsonResponse(w, http.StatusCreated, map[string]interface{}{
		"success": true, "message": "message sent successfully",
	})
}

func (h *ContactHandler) List(w http.ResponseWriter, r *http.Request) {
	messages, err := h.ContactRepo.List(r.Context())
	if err != nil {
		jsonResponse(w, http.StatusInternalServerError, map[string]interface{}{
			"success": false, "message": "failed to list messages",
		})
		return
	}
	unread, _ := h.ContactRepo.CountUnread(r.Context())
	jsonResponse(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"data": map[string]interface{}{
			"messages": messages,
			"unread":   unread,
		},
	})
}

func (h *ContactHandler) MarkRead(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		jsonResponse(w, http.StatusBadRequest, map[string]interface{}{
			"success": false, "message": "invalid id",
		})
		return
	}
	if err := h.ContactRepo.MarkRead(r.Context(), id); err != nil {
		jsonResponse(w, http.StatusInternalServerError, map[string]interface{}{
			"success": false, "message": "failed to mark message as read",
		})
		return
	}
	jsonResponse(w, http.StatusOK, map[string]interface{}{
		"success": true, "message": "message marked as read",
	})
}

func (h *ContactHandler) Delete(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		jsonResponse(w, http.StatusBadRequest, map[string]interface{}{
			"success": false, "message": "invalid id",
		})
		return
	}
	if err := h.ContactRepo.Delete(r.Context(), id); err != nil {
		jsonResponse(w, http.StatusInternalServerError, map[string]interface{}{
			"success": false, "message": "failed to delete message",
		})
		return
	}
	jsonResponse(w, http.StatusOK, map[string]interface{}{
		"success": true, "message": "message deleted",
	})
}
