package handler

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"github.com/aige/api/internal/middleware"
	"github.com/aige/cdk"
	"github.com/aige/order"
	"github.com/aige/payment"
	"github.com/aige/plan"
	"github.com/aige/user"
)

type BillingHandler struct {
	PlanRepo    plan.Repository
	OrderRepo   order.Repository
	UserRepo    user.Repository
	CDKRepo     cdk.Repository
	PaymentRepo payment.Repository
}

func (h *BillingHandler) GetPlans(w http.ResponseWriter, r *http.Request) {
	plans, err := h.PlanRepo.List(r.Context())
	if err != nil {
		jsonResponse(w, http.StatusInternalServerError, map[string]interface{}{
			"success": false, "message": "internal server error",
		})
		return
	}

	jsonResponse(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"data":    plans,
	})
}

func (h *BillingHandler) CreateOrder(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())
	if userID == "" {
		jsonResponse(w, http.StatusUnauthorized, map[string]interface{}{
			"success": false, "message": "unauthorized",
		})
		return
	}

	var req struct {
		PlanID string `json:"plan_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		jsonResponse(w, http.StatusBadRequest, map[string]interface{}{
			"success": false, "message": "invalid request body",
		})
		return
	}

	if req.PlanID == "" {
		jsonResponse(w, http.StatusBadRequest, map[string]interface{}{
			"success": false, "message": "plan_id is required",
		})
		return
	}

	p, err := h.PlanRepo.GetByID(r.Context(), req.PlanID)
	if err != nil {
		jsonResponse(w, http.StatusInternalServerError, map[string]interface{}{
			"success": false, "message": "internal server error",
		})
		return
	}
	if p == nil {
		jsonResponse(w, http.StatusNotFound, map[string]interface{}{
			"success": false, "message": "plan not found",
		})
		return
	}

	orderNo := fmt.Sprintf("%s%d%d", time.Now().Format("20060102150405"), rand.Intn(1000), rand.Intn(1000))

	paymentCfg, err := h.PaymentRepo.GetConfig(r.Context())
	if err != nil || paymentCfg == nil {
		jsonResponse(w, http.StatusInternalServerError, map[string]interface{}{
			"success": false, "message": "payment not configured",
		})
		return
	}

	paymentURL, _, err := payment.EpayCreatePayment(
		paymentCfg.ApiURL, paymentCfg.MerchantID, paymentCfg.MerchantKey,
		paymentCfg.NotifyURL, paymentCfg.ReturnURL,
		orderNo, p.MonthlyPrice, p.Name,
		paymentCfg.Channel,
	)
	if err != nil {
		jsonResponse(w, http.StatusInternalServerError, map[string]interface{}{
			"success": false, "message": "failed to create payment",
		})
		return
	}

	o, err := h.OrderRepo.Create(r.Context(), order.CreateOrderRequest{
		UserID:        userID,
		OrderNo:       orderNo,
		Amount:        p.MonthlyPrice,
		Description:   p.Name,
		PlanID:        p.ID,
		CreditsAmount: p.Credits,
	})
	if err != nil {
		jsonResponse(w, http.StatusInternalServerError, map[string]interface{}{
			"success": false, "message": "failed to create order",
		})
		return
	}

	jsonResponse(w, http.StatusCreated, map[string]interface{}{
		"success": true,
		"data": map[string]interface{}{
			"order":       o,
			"payment_url": paymentURL,
		},
	})
}

func (h *BillingHandler) PaymentNotify(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	orderNo := r.FormValue("out_trade_no")
	_ = r.FormValue("trade_no")
	statusStr := r.FormValue("status")
	sign := r.FormValue("sign")

	if orderNo == "" || sign == "" {
		jsonResponse(w, http.StatusBadRequest, map[string]interface{}{
			"success": false, "message": "invalid notify params",
		})
		return
	}

	o, err := h.OrderRepo.GetByOrderNo(r.Context(), orderNo)
	if err != nil || o == nil {
		w.Write([]byte("fail"))
		return
	}

	paymentCfg, err := h.PaymentRepo.GetConfig(r.Context())
	if err != nil || paymentCfg == nil {
		w.Write([]byte("fail"))
		return
	}

	params := make(map[string]string)
	for k, v := range r.Form {
		if len(v) > 0 {
			params[k] = v[0]
		}
	}

	if !payment.EpayVerifySign(params, paymentCfg.MerchantKey) {
		w.Write([]byte("fail"))
		return
	}

	status, _ := strconv.Atoi(statusStr)
	if status == 1 {
		now := time.Now()
		h.OrderRepo.UpdateStatus(r.Context(), o.ID, order.UpdateOrderStatusRequest{
			Status:        "paid",
			PaymentMethod: paymentCfg.Channel,
			PaymentTime:   &now,
		})

		h.UserRepo.AddCredits(r.Context(), o.UserID, o.CreditsAmount)
	}

	w.Write([]byte("success"))
}

func (h *BillingHandler) RedeemCDK(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())
	if userID == "" {
		jsonResponse(w, http.StatusUnauthorized, map[string]interface{}{
			"success": false, "message": "unauthorized",
		})
		return
	}

	var req cdk.RedeemRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		jsonResponse(w, http.StatusBadRequest, map[string]interface{}{
			"success": false, "message": "invalid request body",
		})
		return
	}

	if req.Code == "" {
		jsonResponse(w, http.StatusBadRequest, map[string]interface{}{
			"success": false, "message": "code is required",
		})
		return
	}

	result, err := h.CDKRepo.Redeem(r.Context(), req.Code, userID)
	if err != nil {
		jsonResponse(w, http.StatusInternalServerError, map[string]interface{}{
			"success": false, "message": err.Error(),
		})
		return
	}

	jsonResponse(w, http.StatusOK, map[string]interface{}{
		"success": result.Credits > 0,
		"data":    result,
	})
}

func (h *BillingHandler) GetOrders(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())
	if userID == "" {
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

	orders, total, err := h.OrderRepo.ListByUserID(r.Context(), userID, offset, limit)
	if err != nil {
		jsonResponse(w, http.StatusInternalServerError, map[string]interface{}{
			"success": false, "message": "internal server error",
		})
		return
	}

	jsonResponse(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"data":    orders,
		"meta": map[string]interface{}{
			"total":  total,
			"offset": offset,
			"limit":  limit,
		},
	})
}

func (h *BillingHandler) GetCredits(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())
	if userID == "" {
		jsonResponse(w, http.StatusUnauthorized, map[string]interface{}{
			"success": false, "message": "unauthorized",
		})
		return
	}

	u, err := h.UserRepo.FindByID(r.Context(), userID)
	if err != nil {
		jsonResponse(w, http.StatusInternalServerError, map[string]interface{}{
			"success": false, "message": "internal server error",
		})
		return
	}
	if u == nil {
		jsonResponse(w, http.StatusNotFound, map[string]interface{}{
			"success": false, "message": "user not found",
		})
		return
	}

	jsonResponse(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"data": map[string]interface{}{
			"credits": u.Credits,
		},
	})
}
