package handler

import (
	"encoding/json"
	"net/http"

	"github.com/aige/setting"
)

type SettingHandler struct {
	SettingRepo setting.Repository
}

func (h *SettingHandler) Get(w http.ResponseWriter, r *http.Request) {
	cfg, err := h.SettingRepo.GetAll(r.Context())
	if err != nil {
		jsonResponse(w, http.StatusInternalServerError, map[string]interface{}{
			"success": false, "message": "failed to get settings",
		})
		return
	}

	jsonResponse(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"data":    cfg,
	})
}

func (h *SettingHandler) Update(w http.ResponseWriter, r *http.Request) {
	var cfg setting.SiteConfig
	if err := json.NewDecoder(r.Body).Decode(&cfg); err != nil {
		jsonResponse(w, http.StatusBadRequest, map[string]interface{}{
			"success": false, "message": "invalid request body",
		})
		return
	}
	if err := setting.ValidateSiteConfig(&cfg); err != nil {
		jsonResponse(w, http.StatusBadRequest, map[string]interface{}{
			"success": false, "message": "站点配置校验失败",
		})
		return
	}

	if err := h.SettingRepo.Update(r.Context(), &cfg); err != nil {
		jsonResponse(w, http.StatusInternalServerError, map[string]interface{}{
			"success": false, "message": "failed to update settings",
		})
		return
	}

	updated, err := h.SettingRepo.GetAll(r.Context())
	if err != nil {
		jsonResponse(w, http.StatusInternalServerError, map[string]interface{}{
			"success": false, "message": "settings updated but failed to read back",
		})
		return
	}

	jsonResponse(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"data":    updated,
	})
}
