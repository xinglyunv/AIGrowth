package project

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// selectCols is the SELECT column list with COALESCE for nullable fields.
const selectCols = `p.id, p.user_id, p.name, COALESCE(p.website, '') AS website, p.industry, COALESCE(p.description, '') AS description, p.keywords, COALESCE(p.service_area, '') AS service_area, COALESCE(p.target_users, '') AS target_users, p.status, p.created_at, p.updated_at,
	bi.brand_intro, bi.product_intro, bi.service_intro, bi.faq::text, bi.advantages::text, bi.cases::text`

const fromClause = `FROM brand_projects p LEFT JOIN brand_infos bi ON bi.project_id = p.id`

// Repository defines the project data access interface.
type Repository interface {
	Create(ctx context.Context, userID string, req CreateProjectRequest) (*BrandProject, error)
	FindByID(ctx context.Context, id string) (*BrandProject, error)
	ListByUser(ctx context.Context, userID string, offset, limit int) ([]BrandProject, int, error)
	Update(ctx context.Context, id string, req UpdateProjectRequest) (*BrandProject, error)
	Delete(ctx context.Context, id string) error
	UpsertBrandInfo(ctx context.Context, projectID string, brandIntro, productIntro, serviceIntro, faq, advantages, cases *string) error
}

// PostgresRepository implements Repository using pgxpool.
type PostgresRepository struct {
	pool *pgxpool.Pool
}

// NewPostgresRepository creates a new PostgresRepository.
func NewPostgresRepository(pool *pgxpool.Pool) Repository {
	return &PostgresRepository{pool: pool}
}

// scanProject scans a row into a BrandProject.
func scanProject(s interface{ Scan(...interface{}) error }) (*BrandProject, error) {
	p := &BrandProject{}
	err := s.Scan(
		&p.ID, &p.UserID, &p.Name, &p.Website, &p.Industry, &p.Description,
		&p.Keywords, &p.ServiceArea, &p.TargetUsers,
		&p.Status, &p.CreatedAt, &p.UpdatedAt,
		&p.BrandIntro, &p.ProductIntro, &p.ServiceIntro,
		&p.FAQ, &p.Advantages, &p.Cases,
	)
	if err != nil {
		return nil, err
	}
	return p, nil
}

// Create inserts a new brand project.
func (r *PostgresRepository) Create(ctx context.Context, userID string, req CreateProjectRequest) (*BrandProject, error) {
	var p BrandProject
	err := r.pool.QueryRow(ctx,
		`INSERT INTO brand_projects (user_id, name, website, industry, description, keywords, service_area, target_users)
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		 RETURNING id, user_id, name, COALESCE(website, '') AS website, industry, COALESCE(description, '') AS description, keywords, COALESCE(service_area, '') AS service_area, COALESCE(target_users, '') AS target_users, status, created_at, updated_at,
		           NULL::text, NULL::text, NULL::text, NULL::text, NULL::text, NULL::text`,
		userID, req.Name, nullStr(req.Website), req.Industry,
		nullStr(req.Description), req.Keywords, nullStr(req.ServiceArea), nullStr(req.TargetUsers),
	).Scan(
		&p.ID, &p.UserID, &p.Name, &p.Website, &p.Industry, &p.Description,
		&p.Keywords, &p.ServiceArea, &p.TargetUsers,
		&p.Status, &p.CreatedAt, &p.UpdatedAt,
		&p.BrandIntro, &p.ProductIntro, &p.ServiceIntro,
		&p.FAQ, &p.Advantages, &p.Cases,
	)
	if err != nil {
		return nil, fmt.Errorf("create project: %w", err)
	}
	return &p, nil
}

// FindByID looks up a project by ID.
func (r *PostgresRepository) FindByID(ctx context.Context, id string) (*BrandProject, error) {
	p, err := scanProject(r.pool.QueryRow(ctx,
		`SELECT `+selectCols+` `+fromClause+` WHERE p.id = $1`, id,
	))
	if err == pgx.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("find project by id: %w", err)
	}
	return p, nil
}

// ListByUser returns paginated projects for a user with a total count.
func (r *PostgresRepository) ListByUser(ctx context.Context, userID string, offset, limit int) ([]BrandProject, int, error) {
	var total int
	err := r.pool.QueryRow(ctx,
		`SELECT COUNT(*) FROM brand_projects WHERE user_id = $1 AND status = 'active'`,
		userID,
	).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("count projects: %w", err)
	}

	rows, err := r.pool.Query(ctx,
		`SELECT `+selectCols+` `+fromClause+`
		 WHERE p.user_id = $1 AND p.status = 'active'
		 ORDER BY p.created_at DESC LIMIT $2 OFFSET $3`,
		userID, limit, offset,
	)
	if err != nil {
		return nil, 0, fmt.Errorf("list projects: %w", err)
	}
	defer rows.Close()

	var projects []BrandProject
	for rows.Next() {
		p, err := scanProject(rows)
		if err != nil {
			return nil, 0, fmt.Errorf("scan project: %w", err)
		}
		projects = append(projects, *p)
	}
	if projects == nil {
		projects = []BrandProject{}
	}
	return projects, total, nil
}

// Update modifies project fields using COALESCE for nullable columns.
func (r *PostgresRepository) Update(ctx context.Context, id string, req UpdateProjectRequest) (*BrandProject, error) {
	tag, err := r.pool.Exec(ctx,
		`UPDATE brand_projects
		 SET name = COALESCE($2, name),
		     website = COALESCE($3, website),
		     industry = COALESCE($4, industry),
		     description = COALESCE($5, description),
		     keywords = COALESCE($6, keywords),
		     service_area = COALESCE($7, service_area),
		     target_users = COALESCE($8, target_users)
		 WHERE id = $1 AND status = 'active'`,
		id,
		req.Name, req.Website, req.Industry, req.Description,
		req.Keywords, req.ServiceArea, req.TargetUsers,
	)
	if err != nil {
		return nil, fmt.Errorf("update project: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return nil, nil
	}
	return r.FindByID(ctx, id)
}

// UpsertBrandInfo inserts or updates brand_infos for a project.
func (r *PostgresRepository) UpsertBrandInfo(ctx context.Context, projectID string, brandIntro, productIntro, serviceIntro, faq, advantages, cases *string) error {
	_, err := r.pool.Exec(ctx,
		`INSERT INTO brand_infos (project_id, brand_intro, product_intro, service_intro, faq, advantages, cases)
		 VALUES ($1, $2, $3, $4, $5::text, $6::text, $7::text)
		 ON CONFLICT (project_id) DO UPDATE SET
		     brand_intro = COALESCE($2, brand_infos.brand_intro),
		     product_intro = COALESCE($3, brand_infos.product_intro),
		     service_intro = COALESCE($4, brand_infos.service_intro),
		     faq = COALESCE($5::text, brand_infos.faq),
		     advantages = COALESCE($6::text, brand_infos.advantages),
		     cases = COALESCE($7::text, brand_infos.cases),
		     updated_at = NOW()`,
		projectID, brandIntro, productIntro, serviceIntro, faq, advantages, cases,
	)
	if err != nil {
		return fmt.Errorf("upsert brand info: %w", err)
	}
	return nil
}

// Delete soft-deletes a project by setting status to 'archived'.
func (r *PostgresRepository) Delete(ctx context.Context, id string) error {
	tag, err := r.pool.Exec(ctx,
		`UPDATE brand_projects SET status = 'archived' WHERE id = $1 AND status = 'active'`,
		id,
	)
	if err != nil {
		return fmt.Errorf("delete project: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return fmt.Errorf("project not found or already archived")
	}
	return nil
}

// nullStr returns nil if the string is empty, otherwise returns the string.
func nullStr(s string) interface{} {
	if s == "" {
		return nil
	}
	return s
}
