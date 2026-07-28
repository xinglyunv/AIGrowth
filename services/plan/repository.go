package plan

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

const planCols = `id, name, code, COALESCE(description, '') AS description, monthly_price, yearly_price, max_projects, max_ai_queries, max_reports, COALESCE(credits, 0) AS credits, COALESCE(features, '') AS features, popular, is_active, sort_order, created_at, updated_at`

type Repository interface {
	List(ctx context.Context) ([]Plan, error)
	GetByID(ctx context.Context, id string) (*Plan, error)
	Create(ctx context.Context, req CreatePlanRequest) (*Plan, error)
	Update(ctx context.Context, id string, req UpdatePlanRequest) (*Plan, error)
}

type PostgresRepository struct {
	pool *pgxpool.Pool
}

func NewPostgresRepository(pool *pgxpool.Pool) Repository {
	return &PostgresRepository{pool: pool}
}

func (r *PostgresRepository) List(ctx context.Context) ([]Plan, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT `+planCols+` FROM plans ORDER BY sort_order ASC`,
	)
	if err != nil {
		return nil, fmt.Errorf("list plans: %w", err)
	}
	defer rows.Close()

	var plans []Plan
	for rows.Next() {
		var p Plan
		err := rows.Scan(
			&p.ID, &p.Name, &p.Code, &p.Description,
			&p.MonthlyPrice, &p.YearlyPrice, &p.MaxProjects, &p.MaxAIQueries,
			&p.MaxReports, &p.Credits, &p.Features, &p.Popular,
			&p.IsActive, &p.SortOrder, &p.CreatedAt, &p.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("scan plan: %w", err)
		}
		plans = append(plans, p)
	}
	if plans == nil {
		plans = []Plan{}
	}
	return plans, nil
}

func (r *PostgresRepository) GetByID(ctx context.Context, id string) (*Plan, error) {
	p := &Plan{}
	err := r.pool.QueryRow(ctx,
		`SELECT `+planCols+` FROM plans WHERE id = $1`, id,
	).Scan(
		&p.ID, &p.Name, &p.Code, &p.Description,
		&p.MonthlyPrice, &p.YearlyPrice, &p.MaxProjects, &p.MaxAIQueries,
		&p.MaxReports, &p.Credits, &p.Features, &p.Popular,
		&p.IsActive, &p.SortOrder, &p.CreatedAt, &p.UpdatedAt,
	)
	if err == pgx.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("get plan by id: %w", err)
	}
	return p, nil
}

func (r *PostgresRepository) Create(ctx context.Context, req CreatePlanRequest) (*Plan, error) {
	p := &Plan{}
	err := r.pool.QueryRow(ctx,
		`INSERT INTO plans (name, code, description, monthly_price, yearly_price, max_projects, max_ai_queries, max_reports, credits, sort_order)
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		 RETURNING `+planCols,
		req.Name, req.Code, req.Description,
		req.MonthlyPrice, req.YearlyPrice,
		req.MaxProjects, req.MaxAIQueries, req.MaxReports, req.Credits, req.SortOrder,
	).Scan(
		&p.ID, &p.Name, &p.Code, &p.Description,
		&p.MonthlyPrice, &p.YearlyPrice, &p.MaxProjects, &p.MaxAIQueries,
		&p.MaxReports, &p.Credits, &p.Features, &p.Popular,
		&p.IsActive, &p.SortOrder, &p.CreatedAt, &p.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("create plan: %w", err)
	}
	return p, nil
}

func (r *PostgresRepository) Update(ctx context.Context, id string, req UpdatePlanRequest) (*Plan, error) {
	setClauses := []string{}
	args := []interface{}{}
	argIdx := 1

	if req.Name != nil {
		setClauses = append(setClauses, fmt.Sprintf("name = $%d", argIdx))
		args = append(args, *req.Name)
		argIdx++
	}
	if req.Code != nil {
		setClauses = append(setClauses, fmt.Sprintf("code = $%d", argIdx))
		args = append(args, *req.Code)
		argIdx++
	}
	if req.Description != nil {
		setClauses = append(setClauses, fmt.Sprintf("description = $%d", argIdx))
		args = append(args, *req.Description)
		argIdx++
	}
	if req.MonthlyPrice != nil {
		setClauses = append(setClauses, fmt.Sprintf("monthly_price = $%d", argIdx))
		args = append(args, *req.MonthlyPrice)
		argIdx++
	}
	if req.YearlyPrice != nil {
		setClauses = append(setClauses, fmt.Sprintf("yearly_price = $%d", argIdx))
		args = append(args, *req.YearlyPrice)
		argIdx++
	}
	if req.MaxProjects != nil {
		setClauses = append(setClauses, fmt.Sprintf("max_projects = $%d", argIdx))
		args = append(args, *req.MaxProjects)
		argIdx++
	}
	if req.MaxAIQueries != nil {
		setClauses = append(setClauses, fmt.Sprintf("max_ai_queries = $%d", argIdx))
		args = append(args, *req.MaxAIQueries)
		argIdx++
	}
	if req.MaxReports != nil {
		setClauses = append(setClauses, fmt.Sprintf("max_reports = $%d", argIdx))
		args = append(args, *req.MaxReports)
		argIdx++
	}
	if req.Credits != nil {
		setClauses = append(setClauses, fmt.Sprintf("credits = $%d", argIdx))
		args = append(args, *req.Credits)
		argIdx++
	}
	if req.Features != nil {
		setClauses = append(setClauses, fmt.Sprintf("features = $%d", argIdx))
		args = append(args, *req.Features)
		argIdx++
	}
	if req.Popular != nil {
		setClauses = append(setClauses, fmt.Sprintf("popular = $%d", argIdx))
		args = append(args, *req.Popular)
		argIdx++
	}
	if req.IsActive != nil {
		setClauses = append(setClauses, fmt.Sprintf("is_active = $%d", argIdx))
		args = append(args, *req.IsActive)
		argIdx++
	}
	if req.SortOrder != nil {
		setClauses = append(setClauses, fmt.Sprintf("sort_order = $%d", argIdx))
		args = append(args, *req.SortOrder)
		argIdx++
	}

	if len(setClauses) == 0 {
		return r.GetByID(ctx, id)
	}

	setClauses = append(setClauses, "updated_at = NOW()")
	args = append(args, id)

	query := fmt.Sprintf(`UPDATE plans SET %s WHERE id = $%d RETURNING `+planCols,
		joinStrings(setClauses, ", "), argIdx)

	p := &Plan{}
	err := r.pool.QueryRow(ctx, query, args...).Scan(
		&p.ID, &p.Name, &p.Code, &p.Description,
		&p.MonthlyPrice, &p.YearlyPrice, &p.MaxProjects, &p.MaxAIQueries,
		&p.MaxReports, &p.Credits, &p.Features, &p.Popular,
		&p.IsActive, &p.SortOrder, &p.CreatedAt, &p.UpdatedAt,
	)
	if err == pgx.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("update plan: %w", err)
	}
	return p, nil
}

func joinStrings(elems []string, sep string) string {
	if len(elems) == 0 {
		return ""
	}
	result := elems[0]
	for _, e := range elems[1:] {
		result += sep + e
	}
	return result
}
