package order

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

const orderCols = `id, user_id, order_no, amount, currency, status, COALESCE(payment_method, '') AS payment_method, payment_time, COALESCE(description, '') AS description, COALESCE(plan_id::text, '') AS plan_id, COALESCE(credits_amount, 0) AS credits_amount, COALESCE(cdk_id::text, '') AS cdk_id, created_at, updated_at`

type Repository interface {
	List(ctx context.Context, offset, limit int) ([]Order, int, error)
	GetByID(ctx context.Context, id string) (*Order, error)
	ListByUserID(ctx context.Context, userID string, offset, limit int) ([]Order, int, error)
	Create(ctx context.Context, req CreateOrderRequest) (*Order, error)
	UpdateStatus(ctx context.Context, id string, req UpdateOrderStatusRequest) (*Order, error)
	GetByOrderNo(ctx context.Context, orderNo string) (*Order, error)
}

type PostgresRepository struct {
	pool *pgxpool.Pool
}

func NewPostgresRepository(pool *pgxpool.Pool) Repository {
	return &PostgresRepository{pool: pool}
}

func (r *PostgresRepository) List(ctx context.Context, offset, limit int) ([]Order, int, error) {
	var total int
	err := r.pool.QueryRow(ctx, `SELECT COUNT(*) FROM orders`).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("count orders: %w", err)
	}

	rows, err := r.pool.Query(ctx,
		`SELECT o.id, o.user_id, o.order_no, o.amount, o.currency, o.status,
		        COALESCE(o.payment_method, '') AS payment_method, o.payment_time,
		        COALESCE(o.description, '') AS description, COALESCE(o.plan_id::text, '') AS plan_id,
		        COALESCE(o.credits_amount, 0) AS credits_amount, COALESCE(o.cdk_id::text, '') AS cdk_id,
		        o.created_at, o.updated_at,
		        COALESCE(u.username, '') AS user_name
		 FROM orders o
		 LEFT JOIN users u ON u.id = o.user_id
		 ORDER BY o.created_at DESC LIMIT $1 OFFSET $2`,
		limit, offset,
	)
	if err != nil {
		return nil, 0, fmt.Errorf("list orders: %w", err)
	}
	defer rows.Close()

	var orders []Order
	for rows.Next() {
		var o Order
		err := rows.Scan(
			&o.ID, &o.UserID, &o.OrderNo, &o.Amount, &o.Currency, &o.Status,
			&o.PaymentMethod, &o.PaymentTime, &o.Description, &o.PlanID,
			&o.CreditsAmount, &o.CDKID,
			&o.CreatedAt, &o.UpdatedAt, &o.UserName,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("scan order: %w", err)
		}
		orders = append(orders, o)
	}
	if orders == nil {
		orders = []Order{}
	}
	return orders, total, nil
}

func (r *PostgresRepository) GetByID(ctx context.Context, id string) (*Order, error) {
	o := &Order{}
	err := r.pool.QueryRow(ctx,
		`SELECT `+orderCols+` FROM orders WHERE id = $1`, id,
	).Scan(
		&o.ID, &o.UserID, &o.OrderNo, &o.Amount, &o.Currency, &o.Status,
		&o.PaymentMethod, &o.PaymentTime, &o.Description, &o.PlanID,
		&o.CreditsAmount, &o.CDKID,
		&o.CreatedAt, &o.UpdatedAt,
	)
	if err == pgx.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("get order by id: %w", err)
	}
	return o, nil
}

func (r *PostgresRepository) ListByUserID(ctx context.Context, userID string, offset, limit int) ([]Order, int, error) {
	var total int
	err := r.pool.QueryRow(ctx,
		`SELECT COUNT(*) FROM orders WHERE user_id = $1`, userID,
	).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("count orders by user: %w", err)
	}

	rows, err := r.pool.Query(ctx,
		`SELECT `+orderCols+` FROM orders WHERE user_id = $1 ORDER BY created_at DESC LIMIT $2 OFFSET $3`,
		userID, limit, offset,
	)
	if err != nil {
		return nil, 0, fmt.Errorf("list orders by user: %w", err)
	}
	defer rows.Close()

	var orders []Order
	for rows.Next() {
		var o Order
		err := rows.Scan(
			&o.ID, &o.UserID, &o.OrderNo, &o.Amount, &o.Currency, &o.Status,
			&o.PaymentMethod, &o.PaymentTime, &o.Description, &o.PlanID,
			&o.CreditsAmount, &o.CDKID,
			&o.CreatedAt, &o.UpdatedAt,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("scan order: %w", err)
		}
		orders = append(orders, o)
	}
	if orders == nil {
		orders = []Order{}
	}
	return orders, total, nil
}

func (r *PostgresRepository) Create(ctx context.Context, req CreateOrderRequest) (*Order, error) {
	currency := req.Currency
	if currency == "" {
		currency = "CNY"
	}
	status := req.Status
	if status == "" {
		status = "pending"
	}

	o := &Order{}
	err := r.pool.QueryRow(ctx,
		`INSERT INTO orders (user_id, order_no, amount, currency, status, description, plan_id, credits_amount, cdk_id)
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		 RETURNING `+orderCols,
		req.UserID, req.OrderNo, req.Amount, currency, status,
		req.Description, nullStr(req.PlanID), req.CreditsAmount, nullStr(req.CDKID),
	).Scan(
		&o.ID, &o.UserID, &o.OrderNo, &o.Amount, &o.Currency, &o.Status,
		&o.PaymentMethod, &o.PaymentTime, &o.Description, &o.PlanID,
		&o.CreditsAmount, &o.CDKID,
		&o.CreatedAt, &o.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("create order: %w", err)
	}
	return o, nil
}

func (r *PostgresRepository) UpdateStatus(ctx context.Context, id string, req UpdateOrderStatusRequest) (*Order, error) {
	o := &Order{}
	err := r.pool.QueryRow(ctx,
		`UPDATE orders SET status = $2, payment_method = COALESCE($3, payment_method), payment_time = COALESCE($4, payment_time), updated_at = NOW()
		 WHERE id = $1
		 RETURNING `+orderCols,
		id, req.Status, req.PaymentMethod, req.PaymentTime,
	).Scan(
		&o.ID, &o.UserID, &o.OrderNo, &o.Amount, &o.Currency, &o.Status,
		&o.PaymentMethod, &o.PaymentTime, &o.Description, &o.PlanID,
		&o.CreditsAmount, &o.CDKID,
		&o.CreatedAt, &o.UpdatedAt,
	)
	if err == pgx.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("update order status: %w", err)
	}
	return o, nil
}

func (r *PostgresRepository) GetByOrderNo(ctx context.Context, orderNo string) (*Order, error) {
	o := &Order{}
	err := r.pool.QueryRow(ctx,
		`SELECT `+orderCols+` FROM orders WHERE order_no = $1`, orderNo,
	).Scan(
		&o.ID, &o.UserID, &o.OrderNo, &o.Amount, &o.Currency, &o.Status,
		&o.PaymentMethod, &o.PaymentTime, &o.Description, &o.PlanID,
		&o.CreditsAmount, &o.CDKID,
		&o.CreatedAt, &o.UpdatedAt,
	)
	if err == pgx.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("get order by order no: %w", err)
	}
	return o, nil
}

func nullStr(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}
