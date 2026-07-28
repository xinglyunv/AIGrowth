package handler

import (
	"net/http"

	"github.com/aige/project"
	"github.com/aige/task"
	"github.com/aige/user"
	"github.com/jackc/pgx/v5/pgxpool"
)

type AdminDashboardHandler struct {
	Pool        *pgxpool.Pool
	UserRepo    user.Repository
	ProjectRepo project.Repository
	TaskRepo    task.Repository
}

type AdminStats struct {
	TotalUsers    int `json:"total_users"`
	TotalProjects int `json:"total_projects"`
	TotalTasks    int `json:"total_tasks"`
}

func (h *AdminDashboardHandler) GetStats(w http.ResponseWriter, r *http.Request) {
	var stats AdminStats

	err := h.Pool.QueryRow(r.Context(), `SELECT COUNT(*) FROM users`).Scan(&stats.TotalUsers)
	if err != nil {
		stats.TotalUsers = 0
	}

	err = h.Pool.QueryRow(r.Context(), `SELECT COUNT(*) FROM brand_projects WHERE status != 'archived'`).Scan(&stats.TotalProjects)
	if err != nil {
		stats.TotalProjects = 0
	}

	err = h.Pool.QueryRow(r.Context(), `SELECT COUNT(*) FROM ai_tasks`).Scan(&stats.TotalTasks)
	if err != nil {
		stats.TotalTasks = 0
	}

	jsonResponse(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"data":    stats,
	})
}
