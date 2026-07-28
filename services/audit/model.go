package audit

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Log struct {
	ID        string                 `json:"id"`
	UserID    string                 `json:"user_id,omitempty"`
	Action    string                 `json:"action"`
	Resource  string                 `json:"resource"`
	Detail    map[string]interface{} `json:"detail,omitempty"`
	IPAddress string                 `json:"ip_address,omitempty"`
	CreatedAt time.Time              `json:"created_at"`
}

type Repository interface {
	Create(ctx context.Context, log *Log) error
	List(ctx context.Context, offset, limit int) ([]*Log, int, error)
}

type PostgresRepository struct {
	pool *pgxpool.Pool
}

func NewPostgresRepository(pool *pgxpool.Pool) Repository {
	return &PostgresRepository{pool: pool}
}

func (r *PostgresRepository) Create(ctx context.Context, log *Log) error {
	_, err := r.pool.Exec(ctx,
		`INSERT INTO audit_logs (user_id, action, resource, detail, ip_address)
		 VALUES ($1, $2, $3, $4, $5)`,
		log.UserID, log.Action, log.Resource, log.Detail, log.IPAddress,
	)
	return err
}

func (r *PostgresRepository) List(ctx context.Context, offset, limit int) ([]*Log, int, error) {
	var total int
	err := r.pool.QueryRow(ctx, `SELECT COUNT(*) FROM audit_logs`).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	rows, err := r.pool.Query(ctx,
		`SELECT id, COALESCE(user_id,''), action, resource, detail, COALESCE(ip_address,''), created_at
		 FROM audit_logs ORDER BY created_at DESC LIMIT $1 OFFSET $2`, limit, offset,
	)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var logs []*Log
	for rows.Next() {
		l := &Log{}
		if err := rows.Scan(&l.ID, &l.UserID, &l.Action, &l.Resource, &l.Detail, &l.IPAddress, &l.CreatedAt); err != nil {
			return nil, 0, err
		}
		logs = append(logs, l)
	}
	return logs, total, rows.Err()
}
