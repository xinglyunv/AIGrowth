package payment

import (
	"context"
	"crypto/md5"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"sort"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository interface {
	GetConfig(ctx context.Context) (*PaymentConfig, error)
	UpdateConfig(ctx context.Context, req UpdatePaymentConfigRequest) (*PaymentConfig, error)
}

type PostgresRepository struct {
	pool *pgxpool.Pool
}

func NewPostgresRepository(pool *pgxpool.Pool) Repository {
	return &PostgresRepository{pool: pool}
}

func (r *PostgresRepository) GetConfig(ctx context.Context) (*PaymentConfig, error) {
	c := &PaymentConfig{}
	err := r.pool.QueryRow(ctx,
		`SELECT id, name, COALESCE(channel, 'alipay'), COALESCE(merchant_id,''), COALESCE(merchant_key,''), COALESCE(api_url,''), COALESCE(notify_url,''), COALESCE(return_url,''), is_active, created_at, updated_at
		 FROM payment_configs ORDER BY created_at DESC LIMIT 1`,
	).Scan(
		&c.ID, &c.Name, &c.Channel, &c.MerchantID, &c.MerchantKey, &c.ApiURL, &c.NotifyURL, &c.ReturnURL,
		&c.IsActive, &c.CreatedAt, &c.UpdatedAt,
	)
	if err == pgx.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("get payment config: %w", err)
	}
	return c, nil
}

func (r *PostgresRepository) UpdateConfig(ctx context.Context, req UpdatePaymentConfigRequest) (*PaymentConfig, error) {
	setClauses := []string{}
	args := []interface{}{}
	argIdx := 1

	if req.MerchantID != nil {
		setClauses = append(setClauses, fmt.Sprintf("merchant_id = $%d", argIdx))
		args = append(args, *req.MerchantID)
		argIdx++
	}
	if req.Channel != nil {
		setClauses = append(setClauses, fmt.Sprintf("channel = $%d", argIdx))
		args = append(args, *req.Channel)
		argIdx++
	}
	if req.MerchantKey != nil {
		setClauses = append(setClauses, fmt.Sprintf("merchant_key = $%d", argIdx))
		args = append(args, *req.MerchantKey)
		argIdx++
	}
	if req.ApiURL != nil {
		setClauses = append(setClauses, fmt.Sprintf("api_url = $%d", argIdx))
		args = append(args, *req.ApiURL)
		argIdx++
	}
	if req.NotifyURL != nil {
		setClauses = append(setClauses, fmt.Sprintf("notify_url = $%d", argIdx))
		args = append(args, *req.NotifyURL)
		argIdx++
	}
	if req.ReturnURL != nil {
		setClauses = append(setClauses, fmt.Sprintf("return_url = $%d", argIdx))
		args = append(args, *req.ReturnURL)
		argIdx++
	}
	if req.IsActive != nil {
		setClauses = append(setClauses, fmt.Sprintf("is_active = $%d", argIdx))
		args = append(args, *req.IsActive)
		argIdx++
	}

	if len(setClauses) == 0 {
		return r.GetConfig(ctx)
	}

	setClauses = append(setClauses, "updated_at = NOW()")

	// Try upsert: insert a default row if none exists, then update
	c := &PaymentConfig{}
	err := r.pool.QueryRow(ctx,
		`INSERT INTO payment_configs (name, channel, merchant_id, merchant_key, api_url, notify_url, return_url, is_active)
		 VALUES ('default', 'alipay', '', '', '', '', '', true)
		 ON CONFLICT DO NOTHING`,
	).Scan(&c.ID, &c.Name)
	if err != nil && err != pgx.ErrNoRows {
		// Ignore "no rows" - the row already exists or insert was skipped
		_ = err
	}

	query := fmt.Sprintf(`UPDATE payment_configs SET %s WHERE id = (SELECT id FROM payment_configs ORDER BY created_at DESC LIMIT 1) RETURNING id, name, COALESCE(channel, 'alipay'), COALESCE(merchant_id,''), COALESCE(merchant_key,''), COALESCE(api_url,''), COALESCE(notify_url,''), COALESCE(return_url,''), is_active, created_at, updated_at`,
		strings.Join(setClauses, ", "))

	c2 := &PaymentConfig{}
	err = r.pool.QueryRow(ctx, query, args...).Scan(
		&c2.ID, &c2.Name, &c2.Channel, &c2.MerchantID, &c2.MerchantKey, &c2.ApiURL, &c2.NotifyURL, &c2.ReturnURL,
		&c2.IsActive, &c2.CreatedAt, &c2.UpdatedAt,
	)
	if err == pgx.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("update payment config: %w", err)
	}
	return c2, nil
}

// EpaySign generates 易支付 sign (MD5)
func EpaySign(params map[string]string, key string) string {
	keys := make([]string, 0, len(params))
	for k := range params {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	var parts []string
	for _, k := range keys {
		if params[k] != "" && k != "sign" && k != "sign_type" {
			parts = append(parts, fmt.Sprintf("%s=%s", k, params[k]))
		}
	}
	parts = append(parts, fmt.Sprintf("key=%s", key))

	data := strings.Join(parts, "&")
	return fmt.Sprintf("%x", md5.Sum([]byte(data)))
}

// EpayCreatePayment calls 易支付 API to create a payment order
// Returns the payment URL for redirect
func EpayCreatePayment(apiURL, merchantID, merchantKey, notifyURL, returnURL, orderNo string, amount float64, description, channel string) (string, string, error) {
	if channel == "" {
		channel = "alipay"
	}
	params := map[string]string{
		"pid":          merchantID,
		"type":         channel,
		"out_trade_no": orderNo,
		"notify_url":   notifyURL,
		"return_url":   returnURL,
		"name":         description,
		"money":        fmt.Sprintf("%.2f", amount),
		"sign_type":    "MD5",
	}

	sign := EpaySign(params, merchantKey)
	params["sign"] = sign

	values := url.Values{}
	for k, v := range params {
		values.Set(k, v)
	}

	paymentURL := fmt.Sprintf("%s/submit.php?%s", strings.TrimRight(apiURL, "/"), values.Encode())

	return paymentURL, "", nil
}

// EpayVerifySign verifies the callback sign from 易支付
func EpayVerifySign(params map[string]string, merchantKey string) bool {
	receivedSign := params["sign"]
	if receivedSign == "" {
		return false
	}
	expectedSign := EpaySign(params, merchantKey)
	return strings.EqualFold(receivedSign, expectedSign)
}

// EpayQueryOrder queries order status from 易支付
func EpayQueryOrder(apiURL, merchantID, merchantKey, orderNo string) (int, error) {
	params := map[string]string{
		"act":          "order",
		"pid":          merchantID,
		"out_trade_no": orderNo,
		"sign_type":    "MD5",
	}
	sign := EpaySign(params, merchantKey)
	params["sign"] = sign

	values := url.Values{}
	for k, v := range params {
		values.Set(k, v)
	}

	queryURL := fmt.Sprintf("%s/api.php?%s", strings.TrimRight(apiURL, "/"), values.Encode())

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Get(queryURL)
	if err != nil {
		return 0, fmt.Errorf("query order: %w", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	var result struct {
		Status int `json:"status"`
	}
	json.Unmarshal(body, &result)

	return result.Status, nil
}
