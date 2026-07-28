package handler

import (
	"log"
	"net/http"
	"strconv"

	"github.com/aige/api/internal/middleware"
	"github.com/aige/task"
	"github.com/go-chi/chi/v5"
)

type AdminTaskHandler struct {
	TaskRepo task.Repository
}

func (h *AdminTaskHandler) ListTasks(w http.ResponseWriter, r *http.Request) {
	adminID := middleware.GetAdminID(r.Context())
	if adminID == "" {
		jsonResponse(w, http.StatusUnauthorized, map[string]interface{}{
			"success": false, "message": "unauthorized",
		})
		return
	}

	offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	if limit <= 0 || limit > 100 {
		limit = 20
	}
	if offset < 0 {
		offset = 0
	}

	tasks, total, err := h.TaskRepo.ListAll(r.Context(), offset, limit)
	if err != nil {
		log.Printf("ERROR listing all tasks: %v", err)
		jsonResponse(w, http.StatusInternalServerError, map[string]interface{}{
			"success": false, "message": "internal server error",
		})
		return
	}

	jsonResponse(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"data":    tasks,
		"meta": map[string]interface{}{
			"total":  total,
			"offset": offset,
			"limit":  limit,
		},
	})
}

func (h *AdminTaskHandler) RetryTask(w http.ResponseWriter, r *http.Request) {
	adminID := middleware.GetAdminID(r.Context())
	if adminID == "" {
		jsonResponse(w, http.StatusUnauthorized, map[string]interface{}{
			"success": false, "message": "unauthorized",
		})
		return
	}

	id := chi.URLParam(r, "id")

	t, err := h.TaskRepo.FindByID(r.Context(), id)
	if err != nil {
		jsonResponse(w, http.StatusInternalServerError, map[string]interface{}{
			"success": false, "message": "internal server error",
		})
		return
	}
	if t == nil {
		jsonResponse(w, http.StatusNotFound, map[string]interface{}{
			"success": false, "message": "task not found",
		})
		return
	}

	if err := h.TaskRepo.UpdateStatus(r.Context(), id, "pending", 0, ""); err != nil {
		log.Printf("ERROR retrying task %s: %v", id, err)
		jsonResponse(w, http.StatusInternalServerError, map[string]interface{}{
			"success": false, "message": "internal server error",
		})
		return
	}

	jsonResponse(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"message": "任务已重置为待执行状态",
	})
}

func (h *AdminTaskHandler) DeleteTask(w http.ResponseWriter, r *http.Request) {
	adminID := middleware.GetAdminID(r.Context())
	if adminID == "" {
		jsonResponse(w, http.StatusUnauthorized, map[string]interface{}{
			"success": false, "message": "unauthorized",
		})
		return
	}

	id := chi.URLParam(r, "id")

	t, err := h.TaskRepo.FindByID(r.Context(), id)
	if err != nil {
		jsonResponse(w, http.StatusInternalServerError, map[string]interface{}{
			"success": false, "message": "internal server error",
		})
		return
	}
	if t == nil {
		jsonResponse(w, http.StatusNotFound, map[string]interface{}{
			"success": false, "message": "task not found",
		})
		return
	}

	if err := h.TaskRepo.Delete(r.Context(), id); err != nil {
		log.Printf("ERROR deleting task %s: %v", id, err)
		jsonResponse(w, http.StatusInternalServerError, map[string]interface{}{
			"success": false, "message": "internal server error",
		})
		return
	}

	jsonResponse(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"message": "任务已删除",
	})
}
