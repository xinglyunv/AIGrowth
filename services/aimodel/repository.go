package aimodel

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository interface {
	List(ctx context.Context) ([]*AIModel, error)
	ListEnabled(ctx context.Context) ([]*AIModel, error)
	GetByID(ctx context.Context, id string) (*AIModel, error)
	Create(ctx context.Context, req *CreateModelRequest) (*AIModel, error)
	Update(ctx context.Context, id string, req *UpdateModelRequest) (*AIModel, error)
	Delete(ctx context.Context, id string) error
	UpdateTestStatus(ctx context.Context, id string, status string) error
}

type PostgresRepository struct {
	pool *pgxpool.Pool
}

func NewPostgresRepository(pool *pgxpool.Pool) Repository {
	return &PostgresRepository{pool: pool}
}

func scanModel(row interface{ Scan(...interface{}) error }) (*AIModel, error) {
	m := &AIModel{}
	err := row.Scan(
		&m.ID, &m.Name, &m.Provider, &m.Model,
		&m.BaseURL, &m.APIKey, &m.Enabled, &m.Description,
		&m.IsSystem, &m.LastTestedAt, &m.LastTestStatus,
		&m.CreatedAt, &m.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return m, nil
}

func (r *PostgresRepository) List(ctx context.Context) ([]*AIModel, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, name, provider, model, base_url, api_key, enabled, description,
		        is_system, last_tested_at, last_test_status, created_at, updated_at
		 FROM ai_models ORDER BY created_at ASC`,
	)
	if err != nil {
		return nil, fmt.Errorf("list ai models: %w", err)
	}
	defer rows.Close()

	models := make([]*AIModel, 0)
	for rows.Next() {
		m, err := scanModel(rows)
		if err != nil {
			return nil, fmt.Errorf("scan ai model: %w", err)
		}
		models = append(models, m)
	}
	return models, rows.Err()
}

func (r *PostgresRepository) ListEnabled(ctx context.Context) ([]*AIModel, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, name, provider, model, base_url, api_key, enabled, description,
		        is_system, last_tested_at, last_test_status, created_at, updated_at
		 FROM ai_models WHERE enabled = TRUE ORDER BY created_at ASC`,
	)
	if err != nil {
		return nil, fmt.Errorf("list enabled models: %w", err)
	}
	defer rows.Close()

	models := make([]*AIModel, 0)
	for rows.Next() {
		m, err := scanModel(rows)
		if err != nil {
			return nil, fmt.Errorf("scan model: %w", err)
		}
		models = append(models, m)
	}
	return models, rows.Err()
}

func (r *PostgresRepository) GetByID(ctx context.Context, id string) (*AIModel, error) {
	m, err := scanModel(r.pool.QueryRow(ctx,
		`SELECT id, name, provider, model, base_url, api_key, enabled, description,
		        is_system, last_tested_at, last_test_status, created_at, updated_at
		 FROM ai_models WHERE id = $1`, id,
	))
	if err != nil {
		return nil, fmt.Errorf("get ai model %s: %w", id, err)
	}
	return m, nil
}

func (r *PostgresRepository) Create(ctx context.Context, req *CreateModelRequest) (*AIModel, error) {
	m, err := scanModel(r.pool.QueryRow(ctx,
		`INSERT INTO ai_models (name, provider, model, base_url, api_key, enabled, description)
		 VALUES ($1, $2, $3, $4, $5, $6, $7)
		 RETURNING id, name, provider, model, base_url, api_key, enabled, description,
		           is_system, last_tested_at, last_test_status, created_at, updated_at`,
		req.Name, req.Provider, req.Model, req.BaseURL, req.APIKey, req.Enabled, req.Description,
	))
	if err != nil {
		return nil, fmt.Errorf("create ai model: %w", err)
	}
	return m, nil
}

func (r *PostgresRepository) Update(ctx context.Context, id string, req *UpdateModelRequest) (*AIModel, error) {
	existing, err := r.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if req.Name != nil { existing.Name = *req.Name }
	if req.Provider != nil { existing.Provider = *req.Provider }
	if req.Model != nil { existing.Model = *req.Model }
	if req.BaseURL != nil { existing.BaseURL = *req.BaseURL }
	if req.APIKey != nil { existing.APIKey = *req.APIKey }
	if req.Enabled != nil { existing.Enabled = *req.Enabled }
	if req.Description != nil { existing.Description = *req.Description }

	existing.UpdatedAt = time.Now()

	m, err := scanModel(r.pool.QueryRow(ctx,
		`UPDATE ai_models SET
			name = $1, provider = $2, model = $3, base_url = $4,
			api_key = $5, enabled = $6, description = $7, updated_at = $8
		 WHERE id = $9
		 RETURNING id, name, provider, model, base_url, api_key, enabled, description,
		           is_system, last_tested_at, last_test_status, created_at, updated_at`,
		existing.Name, existing.Provider, existing.Model,
		existing.BaseURL, existing.APIKey, existing.Enabled, existing.Description,
		existing.UpdatedAt, id,
	))
	if err != nil {
		return nil, fmt.Errorf("update ai model %s: %w", id, err)
	}
	return m, nil
}

func (r *PostgresRepository) Delete(ctx context.Context, id string) error {
	tag, err := r.pool.Exec(ctx, `DELETE FROM ai_models WHERE id = $1`, id)
	if err != nil {
		return fmt.Errorf("delete ai model %s: %w", id, err)
	}
	if tag.RowsAffected() == 0 {
		return fmt.Errorf("ai model %s not found", id)
	}
	return nil
}

func (r *PostgresRepository) UpdateTestStatus(ctx context.Context, id string, status string) error {
	_, err := r.pool.Exec(ctx,
		`UPDATE ai_models SET last_test_status = $1, last_tested_at = NOW(), updated_at = NOW() WHERE id = $2`,
		status, id,
	)
	if err != nil {
		return fmt.Errorf("update test status: %w", err)
	}
	return nil
}
