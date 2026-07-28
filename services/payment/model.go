package payment

import "time"

type PaymentConfig struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Channel     string    `json:"channel"`
	MerchantID  string    `json:"merchant_id"`
	MerchantKey string    `json:"merchant_key"`
	ApiURL      string    `json:"api_url"`
	NotifyURL   string    `json:"notify_url"`
	ReturnURL   string    `json:"return_url"`
	IsActive    bool      `json:"is_active"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type UpdatePaymentConfigRequest struct {
	MerchantID  *string `json:"merchant_id,omitempty"`
	MerchantKey *string `json:"merchant_key,omitempty"`
	Channel     *string `json:"channel,omitempty"`
	ApiURL      *string `json:"api_url,omitempty"`
	NotifyURL   *string `json:"notify_url,omitempty"`
	ReturnURL   *string `json:"return_url,omitempty"`
	IsActive    *bool   `json:"is_active,omitempty"`
}

// 易支付请求参数
type EpayCreateRequest struct {
	PID        int     `json:"pid"`
	Key        string  `json:"key"`
	Type       string  `json:"type"`
	OutTradeNo string  `json:"out_trade_no"`
	NotifyURL  string  `json:"notify_url"`
	ReturnURL  string  `json:"return_url"`
	Name       string  `json:"name"`
	Money      float64 `json:"money"`
	Sign       string  `json:"sign"`
	SignType   string  `json:"sign_type"`
}

// 易支付回调参数
type EpayNotifyParams struct {
	PID        int     `form:"pid"`
	OutTradeNo string  `form:"out_trade_no"`
	TradeNo    string  `form:"trade_no"`
	Money      float64 `form:"money"`
	Status     int     `form:"status"`
	Sign       string  `form:"sign"`
	SignType   string  `form:"sign_type"`
}
