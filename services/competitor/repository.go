package competitor

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

const competitorCols = `id, project_id, name, COALESCE(website, '') AS website, mention_count, rank_position, COALESCE(advantages, '') AS advantages, analysis, created_at, updated_at`

type Repository interface {
	Create(ctx context.Context, comp *Competitor) error
	GetByID(ctx context.Context, id string) (*Competitor, error)
	ListByProjectID(ctx context.Context, projectID string) ([]*Competitor, error)
	Update(ctx context.Context, comp *Competitor) error
	Delete(ctx context.Context, id string) error
	UpsertSummary(ctx context.Context, projectID, name string, mentionCount, rankPosition int, advantages string, analysis map[string]interface{}) error
}

type PostgresRepository struct {
	pool *pgxpool.Pool
}

func NewPostgresRepository(pool *pgxpool.Pool) Repository {
	return &PostgresRepository{pool: pool}
}

func scanCompetitor(row interface{ Scan(...interface{}) error }) (*Competitor, error) {
	c := &Competitor{}
	err := row.Scan(
		&c.ID, &c.ProjectID, &c.Name, &c.Website,
		&c.MentionCount, &c.RankPosition, &c.Advantages,
		&c.Analysis, &c.CreatedAt, &c.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return c, nil
}

func (r *PostgresRepository) Create(ctx context.Context, comp *Competitor) error {
	_, err := scanCompetitor(r.pool.QueryRow(ctx,
		`INSERT INTO competitors (project_id, name, website, mention_count, rank_position, advantages, analysis)
		 VALUES ($1, $2, $3, $4, $5, $6, $7)
		 RETURNING `+competitorCols,
		comp.ProjectID, comp.Name, nullStr(comp.Website),
		comp.MentionCount, comp.RankPosition, nullStr(comp.Advantages), comp.Analysis,
	))
	if err != nil {
		return fmt.Errorf("create competitor: %w", err)
	}
	return nil
}

func (r *PostgresRepository) GetByID(ctx context.Context, id string) (*Competitor, error) {
	c, err := scanCompetitor(r.pool.QueryRow(ctx,
		`SELECT `+competitorCols+` FROM competitors WHERE id = $1`, id,
	))
	if err == pgx.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("get competitor %s: %w", id, err)
	}
	return c, nil
}

func (r *PostgresRepository) ListByProjectID(ctx context.Context, projectID string) ([]*Competitor, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT `+competitorCols+` FROM competitors WHERE project_id = $1 ORDER BY mention_count DESC`,
		projectID,
	)
	if err != nil {
		return nil, fmt.Errorf("list competitors by project: %w", err)
	}
	defer rows.Close()

	competitors := make([]*Competitor, 0)
	for rows.Next() {
		c, err := scanCompetitor(rows)
		if err != nil {
			return nil, fmt.Errorf("scan competitor: %w", err)
		}
		competitors = append(competitors, c)
	}
	return competitors, rows.Err()
}

func (r *PostgresRepository) Update(ctx context.Context, comp *Competitor) error {
	_, err := scanCompetitor(r.pool.QueryRow(ctx,
		`UPDATE competitors SET name = $2, website = $3, mention_count = $4,
		       rank_position = $5, advantages = $6, analysis = $7, updated_at = NOW()
		 WHERE id = $1
		 RETURNING `+competitorCols,
		comp.ID, comp.Name, nullStr(comp.Website),
		comp.MentionCount, comp.RankPosition, nullStr(comp.Advantages), comp.Analysis,
	))
	if err != nil {
		return fmt.Errorf("update competitor %s: %w", comp.ID, err)
	}
	return nil
}

func (r *PostgresRepository) Delete(ctx context.Context, id string) error {
	tag, err := r.pool.Exec(ctx, `DELETE FROM competitors WHERE id = $1`, id)
	if err != nil {
		return fmt.Errorf("delete competitor %s: %w", id, err)
	}
	if tag.RowsAffected() == 0 {
		return fmt.Errorf("competitor %s not found", id)
	}
	return nil
}

func (r *PostgresRepository) UpsertSummary(ctx context.Context, projectID, name string, mentionCount, rankPosition int, advantages string, analysis map[string]interface{}) error {
	_, err := r.pool.Exec(ctx,
		`INSERT INTO competitors (project_id, name, mention_count, rank_position, advantages, analysis)
		 VALUES ($1, $2, $3, NULLIF($4, 0), $5, $6)
		 ON CONFLICT (project_id, name) DO UPDATE SET
		   mention_count = competitors.mention_count + EXCLUDED.mention_count,
		   rank_position = CASE
		     WHEN competitors.rank_position IS NULL OR EXCLUDED.rank_position IS NULL THEN COALESCE(competitors.rank_position, EXCLUDED.rank_position)
		     WHEN EXCLUDED.rank_position < competitors.rank_position THEN EXCLUDED.rank_position
		     ELSE competitors.rank_position
		   END,
		   advantages = COALESCE(NULLIF(EXCLUDED.advantages, ''), competitors.advantages),
		   analysis = COALESCE(competitors.analysis, '{}'::jsonb) || COALESCE(EXCLUDED.analysis, '{}'::jsonb),
		   updated_at = NOW()`,
		projectID, name, mentionCount, rankPosition, nullStr(advantages), analysis)
	if err != nil {
		return fmt.Errorf("upsert competitor summary: %w", err)
	}
	return nil
}

func nullStr(s string) interface{} {
	if s == "" {
		return nil
	}
	return s
}
