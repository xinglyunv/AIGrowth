package report

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

const reportCols = `id, project_id, task_id, user_id, title, type, visibility_score, content, COALESCE(summary, '') AS summary, status, COALESCE(share_token, '') AS share_token, share_expires_at, created_at, updated_at`

type Repository interface {
	Create(ctx context.Context, r *Report) error
	GetByID(ctx context.Context, id string) (*Report, error)
	ListByUserID(ctx context.Context, userID string, offset, limit int) ([]*Report, int, error)
	ListByProjectID(ctx context.Context, projectID string, offset, limit int) ([]*Report, int, error)
	UpdateStatus(ctx context.Context, id, status string) error
	UpdateShareToken(ctx context.Context, id, shareToken string, shareExpiresAt *time.Time) error
}

type PostgresRepository struct {
	pool *pgxpool.Pool
}

func NewPostgresRepository(pool *pgxpool.Pool) Repository {
	return &PostgresRepository{pool: pool}
}

func scanReport(row interface{ Scan(...interface{}) error }) (*Report, error) {
	r := &Report{}
	err := row.Scan(
		&r.ID, &r.ProjectID, &r.TaskID, &r.UserID,
		&r.Title, &r.Type, &r.VisibilityScore, &r.Content,
		&r.Summary, &r.Status, &r.ShareToken, &r.ShareExpiresAt,
		&r.CreatedAt, &r.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return r, nil
}

func (r *PostgresRepository) Create(ctx context.Context, report *Report) error {
	_, err := scanReport(r.pool.QueryRow(ctx,
		`INSERT INTO reports (project_id, task_id, user_id, title, type, visibility_score, content, summary, status)
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		 RETURNING `+reportCols,
		report.ProjectID, report.TaskID, report.UserID,
		report.Title, report.Type, report.VisibilityScore,
		report.Content, nullStr(report.Summary), report.Status,
	))
	if err != nil {
		return fmt.Errorf("create report: %w", err)
	}
	return nil
}

func (r *PostgresRepository) GetByID(ctx context.Context, id string) (*Report, error) {
	report, err := scanReport(r.pool.QueryRow(ctx,
		`SELECT `+reportCols+` FROM reports WHERE id = $1`, id,
	))
	if err == pgx.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("get report %s: %w", id, err)
	}
	return report, nil
}

func (r *PostgresRepository) ListByUserID(ctx context.Context, userID string, offset, limit int) ([]*Report, int, error) {
	var total int
	err := r.pool.QueryRow(ctx,
		`SELECT COUNT(*) FROM reports WHERE user_id = $1`, userID,
	).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("count reports by user: %w", err)
	}

	rows, err := r.pool.Query(ctx,
		`SELECT `+reportCols+` FROM reports
		 WHERE user_id = $1
		 ORDER BY created_at DESC LIMIT $2 OFFSET $3`,
		userID, limit, offset,
	)
	if err != nil {
		return nil, 0, fmt.Errorf("list reports by user: %w", err)
	}
	defer rows.Close()

	reports := make([]*Report, 0)
	for rows.Next() {
		rpt, err := scanReport(rows)
		if err != nil {
			return nil, 0, fmt.Errorf("scan report: %w", err)
		}
		reports = append(reports, rpt)
	}
	return reports, total, rows.Err()
}

func (r *PostgresRepository) ListByProjectID(ctx context.Context, projectID string, offset, limit int) ([]*Report, int, error) {
	var total int
	err := r.pool.QueryRow(ctx,
		`SELECT COUNT(*) FROM reports WHERE project_id = $1`, projectID,
	).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("count reports by project: %w", err)
	}

	rows, err := r.pool.Query(ctx,
		`SELECT `+reportCols+` FROM reports
		 WHERE project_id = $1
		 ORDER BY created_at DESC LIMIT $2 OFFSET $3`,
		projectID, limit, offset,
	)
	if err != nil {
		return nil, 0, fmt.Errorf("list reports by project: %w", err)
	}
	defer rows.Close()

	reports := make([]*Report, 0)
	for rows.Next() {
		rpt, err := scanReport(rows)
		if err != nil {
			return nil, 0, fmt.Errorf("scan report: %w", err)
		}
		reports = append(reports, rpt)
	}
	return reports, total, rows.Err()
}

func (r *PostgresRepository) UpdateStatus(ctx context.Context, id, status string) error {
	tag, err := r.pool.Exec(ctx,
		`UPDATE reports SET status = $2, updated_at = NOW() WHERE id = $1`,
		id, status,
	)
	if err != nil {
		return fmt.Errorf("update report status: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return fmt.Errorf("report %s not found", id)
	}
	return nil
}

func (r *PostgresRepository) UpdateShareToken(ctx context.Context, id, shareToken string, shareExpiresAt *time.Time) error {
	tag, err := r.pool.Exec(ctx,
		`UPDATE reports SET share_token = $2, share_expires_at = $3, updated_at = NOW() WHERE id = $1`,
		id, shareToken, shareExpiresAt,
	)
	if err != nil {
		return fmt.Errorf("update report share token: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return fmt.Errorf("report %s not found", id)
	}
	return nil
}

func nullStr(s string) interface{} {
	if s == "" {
		return nil
	}
	return s
}
