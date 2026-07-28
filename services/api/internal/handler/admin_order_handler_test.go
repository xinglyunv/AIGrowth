package handler

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/aige/api/internal/middleware"
	"github.com/aige/order"
)

type stubOrderRepository struct {
	orders []order.Order
	total  int
	err    error
}

func (s stubOrderRepository) List(context.Context, int, int) ([]order.Order, int, error) {
	return s.orders, s.total, s.err
}

func (stubOrderRepository) GetByID(context.Context, string) (*order.Order, error) {
	return nil, nil
}

func (stubOrderRepository) ListByUserID(context.Context, string, int, int) ([]order.Order, int, error) {
	return nil, 0, nil
}

func (stubOrderRepository) Create(context.Context, order.CreateOrderRequest) (*order.Order, error) {
	return nil, nil
}

func (stubOrderRepository) UpdateStatus(context.Context, string, order.UpdateOrderStatusRequest) (*order.Order, error) {
	return nil, nil
}

func (stubOrderRepository) GetByOrderNo(context.Context, string) (*order.Order, error) {
	return nil, nil
}

func withAdminID(req *http.Request) *http.Request {
	return req.WithContext(context.WithValue(req.Context(), middleware.AdminIDKey, "admin-id"))
}

func TestAdminOrderHandlerListOrdersReturnsPaginatedData(t *testing.T) {
	h := &AdminOrderHandler{OrderRepo: stubOrderRepository{
		orders: []order.Order{{ID: "order-id", OrderNo: "AIGE-1", Status: "paid"}},
		total:  1,
	}}
	req := withAdminID(httptest.NewRequest(http.MethodGet, "/orders?offset=0&limit=20", nil))
	rec := httptest.NewRecorder()

	h.ListOrders(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rec.Code)
	}
	var body struct {
		Success bool          `json:"success"`
		Data    []order.Order `json:"data"`
		Meta    struct {
			Total int `json:"total"`
		} `json:"meta"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if !body.Success || len(body.Data) != 1 || body.Meta.Total != 1 {
		t.Fatalf("unexpected response: %+v", body)
	}
}

func TestAdminOrderHandlerListOrdersHidesRepositoryError(t *testing.T) {
	h := &AdminOrderHandler{OrderRepo: stubOrderRepository{err: errors.New("database detail")}}
	req := withAdminID(httptest.NewRequest(http.MethodGet, "/orders", nil))
	rec := httptest.NewRecorder()

	h.ListOrders(rec, req)

	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("expected status 500, got %d", rec.Code)
	}
	if string(rec.Body.Bytes()) == "" || string(rec.Body.Bytes()) == "database detail" {
		t.Fatalf("repository error leaked in response: %s", rec.Body.String())
	}
}

func TestStatsIntervalNormalizesSupportedRanges(t *testing.T) {
	tests := map[string]string{
		"today": "1 day",
		"7d":    "7 days",
		"30d":   "30 days",
		"":      "7 days",
		"bad":   "7 days",
	}
	for input, expected := range tests {
		if actual := statsInterval(input); actual != expected {
			t.Errorf("statsInterval(%q) = %q, want %q", input, actual, expected)
		}
	}
}
