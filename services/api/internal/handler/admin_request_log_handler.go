package handler

import (
	"net/http"
	"strconv"

	"github.com/aige/requestlog"
)

type AdminRequestLogHandler struct {
	RequestLogRepo requestlog.Repository
}

func NewAdminRequestLogHandler(repo requestlog.Repository) *AdminRequestLogHandler {
	return &AdminRequestLogHandler{RequestLogRepo: repo}
}

func (h *AdminRequestLogHandler) List(w http.ResponseWriter, r *http.Request) {
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	if page < 1 {
		page = 1
	}
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	if limit < 1 || limit > 100 {
		limit = 20
	}
	offset := (page - 1) * limit

	logs, total, err := h.RequestLogRepo.List(r.Context(), offset, limit)
	if err != nil {
		jsonResponse(w, http.StatusInternalServerError, map[string]interface{}{"success": false, "message": "failed to list request logs"})
		return
	}

	jsonResponse(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"data": map[string]interface{}{
			"logs":  logs,
			"total": total,
			"page":  page,
			"limit": limit,
		},
	})
}
