package order

import "time"

type Order struct {
	ID            string     `json:"id"`
	UserID        string     `json:"user_id"`
	OrderNo       string     `json:"order_no"`
	Amount        float64    `json:"amount"`
	Currency      string     `json:"currency"`
	Status        string     `json:"status"`
	PaymentMethod string     `json:"payment_method,omitempty"`
	PaymentTime   *time.Time `json:"payment_time,omitempty"`
	Description   string     `json:"description,omitempty"`
	PlanID        string     `json:"plan_id,omitempty"`
	CreditsAmount int        `json:"credits_amount"`
	CDKID         string     `json:"cdk_id,omitempty"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
	UserName      string     `json:"user_name,omitempty"`
}

type CreateOrderRequest struct {
	UserID      string  `json:"user_id"`
	OrderNo     string  `json:"order_no"`
	Amount      float64 `json:"amount"`
	Currency    string  `json:"currency,omitempty"`
	Status      string  `json:"status,omitempty"`
	Description string  `json:"description,omitempty"`
	PlanID        string  `json:"plan_id,omitempty"`
	CreditsAmount int     `json:"credits_amount,omitempty"`
	CDKID         string  `json:"cdk_id,omitempty"`
}

type UpdateOrderStatusRequest struct {
	Status        string     `json:"status"`
	PaymentMethod string     `json:"payment_method,omitempty"`
	PaymentTime   *time.Time `json:"payment_time,omitempty"`
}
