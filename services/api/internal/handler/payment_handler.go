package handler

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/aige/api/internal/middleware"
	"github.com/aige/payment"
)

type AdminPaymentHandler struct {
	PaymentRepo payment.Repository
}

func (h *AdminPaymentHandler) GetConfig(w http.ResponseWriter, r *http.Request) {
	adminID := middleware.GetAdminID(r.Context())
	if adminID == "" {
		jsonResponse(w, http.StatusUnauthorized, map[string]interface{}{
			"success": false, "message": "unauthorized",
		})
		return
	}

	cfg, err := h.PaymentRepo.GetConfig(r.Context())
	if err != nil {
		log.Printf("[PAYMENT] GetConfig error: %v", err)
		jsonResponse(w, http.StatusInternalServerError, map[string]interface{}{
			"success": false, "message": fmt.Sprintf("获取支付配置失败: %v", err),
		})
		return
	}
	if cfg == nil {
		cfg = &payment.PaymentConfig{}
	}

	jsonResponse(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"data":    cfg,
	})
}

func (h *AdminPaymentHandler) UpdateConfig(w http.ResponseWriter, r *http.Request) {
	adminID := middleware.GetAdminID(r.Context())
	if adminID == "" {
		jsonResponse(w, http.StatusUnauthorized, map[string]interface{}{
			"success": false, "message": "unauthorized",
		})
		return
	}

	var req payment.UpdatePaymentConfigRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		jsonResponse(w, http.StatusBadRequest, map[string]interface{}{
			"success": false, "message": "invalid request body",
		})
		return
	}

	cfg, err := h.PaymentRepo.UpdateConfig(r.Context(), req)
	if err != nil {
		log.Printf("[PAYMENT] UpdateConfig error: %v", err)
		jsonResponse(w, http.StatusInternalServerError, map[string]interface{}{
			"success": false, "message": fmt.Sprintf("保存支付配置失败: %v", err),
		})
		return
	}
	if cfg == nil {
		cfg = &payment.PaymentConfig{}
	}

	jsonResponse(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"data":    cfg,
	})
}
