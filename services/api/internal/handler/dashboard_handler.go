package handler

import (
	"net/http"

	"github.com/aige/api/internal/middleware"
	"github.com/aige/project"
	"github.com/aige/task"
	"github.com/aige/user"
)

type DashboardHandler struct {
	TaskRepo    task.Repository
	UserRepo    user.Repository
	ProjectRepo project.Repository
}

func (h *DashboardHandler) GetStats(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())
	if userID == "" {
		jsonResponse(w, http.StatusUnauthorized, map[string]interface{}{
			"success": false, "message": "unauthorized",
		})
		return
	}

	stats, err := h.TaskRepo.GetDashboardStats(r.Context(), userID)
	if err != nil {
		jsonResponse(w, http.StatusInternalServerError, map[string]interface{}{
			"success": false, "message": "internal server error",
		})
		return
	}

	jsonResponse(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"data":    stats,
	})
}
