package requestlog

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository interface {
	Insert(ctx context.Context, log *RequestLog) error
	List(ctx context.Context, offset, limit int) ([]RequestLog, int, error)
}

type PostgresRepository struct {
	pool *pgxpool.Pool
}

func NewPostgresRepository(pool *pgxpool.Pool) Repository {
	return &PostgresRepository{pool: pool}
}

func (r *PostgresRepository) Insert(ctx context.Context, log *RequestLog) error {
	_, err := r.pool.Exec(ctx,
		`INSERT INTO request_logs (method, path, status_code, duration_ms, request_body, response_body, admin_id, user_id, ip_address)
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`,
		log.Method, log.Path, log.StatusCode, log.DurationMs, log.RequestBody, log.ResponseBody, log.AdminID, log.UserID, log.IPAddress)
	return err
}

func (r *PostgresRepository) List(ctx context.Context, offset, limit int) ([]RequestLog, int, error) {
	var total int
	r.pool.QueryRow(ctx, `SELECT COUNT(*) FROM request_logs`).Scan(&total)

	rows, err := r.pool.Query(ctx,
		`SELECT id, method, path, status_code, duration_ms, COALESCE(request_body,''), COALESCE(response_body,''), COALESCE(admin_id::text,''), COALESCE(user_id::text,''), COALESCE(ip_address,''), created_at
		 FROM request_logs ORDER BY created_at DESC LIMIT $1 OFFSET $2`, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var logs []RequestLog
	for rows.Next() {
		var l RequestLog
		if err := rows.Scan(&l.ID, &l.Method, &l.Path, &l.StatusCode, &l.DurationMs, &l.RequestBody, &l.ResponseBody, &l.AdminID, &l.UserID, &l.IPAddress, &l.CreatedAt); err != nil {
			return nil, 0, err
		}
		logs = append(logs, l)
	}
	if logs == nil {
		logs = []RequestLog{}
	}
	return logs, total, rows.Err()
}
