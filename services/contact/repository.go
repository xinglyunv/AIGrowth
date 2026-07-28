package contact

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository interface {
	Create(ctx context.Context, req *CreateMessageRequest) error
	List(ctx context.Context) ([]*Message, error)
	MarkRead(ctx context.Context, id int64) error
	Delete(ctx context.Context, id int64) error
	CountUnread(ctx context.Context) (int, error)
}

type PostgresRepository struct {
	pool *pgxpool.Pool
}

func NewPostgresRepository(pool *pgxpool.Pool) Repository {
	return &PostgresRepository{pool: pool}
}

func (r *PostgresRepository) Create(ctx context.Context, req *CreateMessageRequest) error {
	_, err := r.pool.Exec(ctx,
		`INSERT INTO contact_messages (name, email, subject, message) VALUES ($1, $2, $3, $4)`,
		req.Name, req.Email, req.Subject, req.Message,
	)
	if err != nil {
		return fmt.Errorf("create contact message: %w", err)
	}
	return nil
}

func (r *PostgresRepository) List(ctx context.Context) ([]*Message, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, name, email, subject, message, is_read, created_at, updated_at
		 FROM contact_messages ORDER BY created_at DESC`,
	)
	if err != nil {
		return nil, fmt.Errorf("list contact messages: %w", err)
	}
	defer rows.Close()

	messages := make([]*Message, 0)
	for rows.Next() {
		m := &Message{}
		if err := rows.Scan(&m.ID, &m.Name, &m.Email, &m.Subject, &m.Message, &m.IsRead, &m.CreatedAt, &m.UpdatedAt); err != nil {
			return nil, fmt.Errorf("scan contact message: %w", err)
		}
		messages = append(messages, m)
	}
	return messages, rows.Err()
}

func (r *PostgresRepository) MarkRead(ctx context.Context, id int64) error {
	_, err := r.pool.Exec(ctx,
		`UPDATE contact_messages SET is_read = TRUE, updated_at = NOW() WHERE id = $1`, id,
	)
	if err != nil {
		return fmt.Errorf("mark message read: %w", err)
	}
	return nil
}

func (r *PostgresRepository) Delete(ctx context.Context, id int64) error {
	_, err := r.pool.Exec(ctx, `DELETE FROM contact_messages WHERE id = $1`, id)
	if err != nil {
		return fmt.Errorf("delete contact message: %w", err)
	}
	return nil
}

func (r *PostgresRepository) CountUnread(ctx context.Context) (int, error) {
	var count int
	err := r.pool.QueryRow(ctx, `SELECT COUNT(*) FROM contact_messages WHERE is_read = FALSE`).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("count unread messages: %w", err)
	}
	return count, nil
}
