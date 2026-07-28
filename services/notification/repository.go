package notification

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

const notifCols = `id, user_id, type, title, content, is_read, related_id, created_at`

type Repository interface {
	Create(ctx context.Context, n *Notification) error
	ListByUserID(ctx context.Context, userID string, offset, limit int) ([]*Notification, int, error)
	MarkRead(ctx context.Context, id string) error
	MarkAllRead(ctx context.Context, userID string) error
	CountUnread(ctx context.Context, userID string) (int, error)
}

type PostgresRepository struct {
	pool *pgxpool.Pool
}

func NewPostgresRepository(pool *pgxpool.Pool) Repository {
	return &PostgresRepository{pool: pool}
}

func scanNotification(row interface{ Scan(...interface{}) error }) (*Notification, error) {
	n := &Notification{}
	var relatedID *string
	err := row.Scan(
		&n.ID, &n.UserID, &n.Type, &n.Title,
		&n.Content, &n.IsRead, &relatedID, &n.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	if relatedID != nil {
		n.RelatedID = *relatedID
	}
	return n, nil
}

func (r *PostgresRepository) Create(ctx context.Context, n *Notification) error {
	var relatedID interface{}
	if n.RelatedID != "" {
		relatedID = n.RelatedID
	}
	_, err := scanNotification(r.pool.QueryRow(ctx,
		`INSERT INTO notifications (user_id, type, title, content, related_id)
		 VALUES ($1, $2, $3, $4, $5)
		 RETURNING `+notifCols,
		n.UserID, n.Type, n.Title, n.Content, relatedID,
	))
	if err != nil {
		return fmt.Errorf("create notification: %w", err)
	}
	return nil
}

func (r *PostgresRepository) ListByUserID(ctx context.Context, userID string, offset, limit int) ([]*Notification, int, error) {
	var total int
	err := r.pool.QueryRow(ctx,
		`SELECT COUNT(*) FROM notifications WHERE user_id = $1`, userID,
	).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("count notifications: %w", err)
	}

	rows, err := r.pool.Query(ctx,
		`SELECT `+notifCols+` FROM notifications
		 WHERE user_id = $1
		 ORDER BY created_at DESC LIMIT $2 OFFSET $3`,
		userID, limit, offset,
	)
	if err != nil {
		return nil, 0, fmt.Errorf("list notifications: %w", err)
	}
	defer rows.Close()

	notifs := make([]*Notification, 0)
	for rows.Next() {
		n, err := scanNotification(rows)
		if err != nil {
			return nil, 0, fmt.Errorf("scan notification: %w", err)
		}
		notifs = append(notifs, n)
	}
	return notifs, total, rows.Err()
}

func (r *PostgresRepository) MarkRead(ctx context.Context, id string) error {
	tag, err := r.pool.Exec(ctx,
		`UPDATE notifications SET is_read = TRUE WHERE id = $1`, id,
	)
	if err != nil {
		return fmt.Errorf("mark notification read: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return fmt.Errorf("notification %s not found", id)
	}
	return nil
}

func (r *PostgresRepository) MarkAllRead(ctx context.Context, userID string) error {
	_, err := r.pool.Exec(ctx,
		`UPDATE notifications SET is_read = TRUE WHERE user_id = $1 AND is_read = FALSE`,
		userID,
	)
	if err != nil {
		return fmt.Errorf("mark all notifications read: %w", err)
	}
	return nil
}

func (r *PostgresRepository) CountUnread(ctx context.Context, userID string) (int, error) {
	var count int
	err := r.pool.QueryRow(ctx,
		`SELECT COUNT(*) FROM notifications WHERE user_id = $1 AND is_read = FALSE`,
		userID,
	).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("count unread notifications: %w", err)
	}
	return count, nil
}

func nullStr(s string) interface{} {
	if s == "" {
		return nil
	}
	return s
}
