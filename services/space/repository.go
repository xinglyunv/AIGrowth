package space

import (
	"context"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository interface {
	ListByUser(context.Context, string) ([]Space, error)
	Create(context.Context, string, CreateRequest) (*Space, error)
	GetCurrent(context.Context, string) (*Space, error)
	SetCurrent(context.Context, string, string) (*Space, error)
	ListMembers(context.Context, string, string) ([]Member, error)
	Invite(context.Context, string, string, InviteRequest) (bool, error)
	UpdateMemberRole(context.Context, string, string, string, string) (bool, error)
	RemoveMember(context.Context, string, string, string) (bool, error)
}

type PostgresRepository struct{ pool *pgxpool.Pool }

func NewPostgresRepository(pool *pgxpool.Pool) Repository { return &PostgresRepository{pool: pool} }

func (r *PostgresRepository) ListByUser(ctx context.Context, userID string) ([]Space, error) {
	rows, err := r.pool.Query(ctx, `SELECT s.id, s.owner_id, s.name, s.slug, s.status, s.created_at, s.updated_at
		FROM spaces s JOIN space_members sm ON sm.space_id = s.id
		WHERE sm.user_id = $1 AND sm.status = 'active' AND s.status = 'active' ORDER BY s.created_at`, userID)
	if err != nil {
		return nil, fmt.Errorf("list spaces: %w", err)
	}
	defer rows.Close()
	spaces := make([]Space, 0)
	for rows.Next() {
		var s Space
		if err := rows.Scan(&s.ID, &s.OwnerID, &s.Name, &s.Slug, &s.Status, &s.CreatedAt, &s.UpdatedAt); err != nil {
			return nil, fmt.Errorf("scan space: %w", err)
		}
		spaces = append(spaces, s)
	}
	return spaces, rows.Err()
}

func (r *PostgresRepository) Create(ctx context.Context, userID string, req CreateRequest) (*Space, error) {
	slug := slugify(req.Name)
	if len(userID) >= 8 {
		slug += "-" + strings.ReplaceAll(userID[:8], "-", "")
	}
	var s Space
	err := r.pool.QueryRow(ctx, `WITH created AS (
		INSERT INTO spaces (owner_id, name, slug) VALUES ($1, $2, $3)
		RETURNING id, owner_id, name, slug, status, created_at, updated_at
	) INSERT INTO space_members (space_id, user_id, role) SELECT id, $1, 'owner' FROM created
	RETURNING (SELECT id FROM created), (SELECT owner_id FROM created), (SELECT name FROM created), (SELECT slug FROM created), (SELECT status FROM created), (SELECT created_at FROM created), (SELECT updated_at FROM created)`, userID, req.Name, slug).
		Scan(&s.ID, &s.OwnerID, &s.Name, &s.Slug, &s.Status, &s.CreatedAt, &s.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("create space: %w", err)
	}
	return &s, nil
}

func (r *PostgresRepository) GetCurrent(ctx context.Context, userID string) (*Space, error) {
	var s Space
	err := r.pool.QueryRow(ctx, `SELECT s.id, s.owner_id, s.name, s.slug, s.status, s.created_at, s.updated_at
		FROM users u JOIN spaces s ON s.id = u.current_space_id
		JOIN space_members sm ON sm.space_id = s.id AND sm.user_id = u.id
		WHERE u.id = $1 AND s.status = 'active'`, userID).
		Scan(&s.ID, &s.OwnerID, &s.Name, &s.Slug, &s.Status, &s.CreatedAt, &s.UpdatedAt)
	if err == pgx.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("get current space: %w", err)
	}
	return &s, nil
}

func (r *PostgresRepository) SetCurrent(ctx context.Context, userID, spaceID string) (*Space, error) {
	var s Space
	err := r.pool.QueryRow(ctx, `WITH selected AS (
		SELECT s.id, s.owner_id, s.name, s.slug, s.status, s.created_at, s.updated_at
		FROM spaces s JOIN space_members sm ON sm.space_id = s.id
		WHERE s.id = $2 AND sm.user_id = $1 AND sm.status = 'active' AND s.status = 'active'
	) UPDATE users u SET current_space_id = (SELECT id FROM selected) WHERE u.id = $1 AND EXISTS (SELECT 1 FROM selected)
	RETURNING (SELECT id FROM selected), (SELECT owner_id FROM selected), (SELECT name FROM selected), (SELECT slug FROM selected), (SELECT status FROM selected), (SELECT created_at FROM selected), (SELECT updated_at FROM selected)`, userID, spaceID).
		Scan(&s.ID, &s.OwnerID, &s.Name, &s.Slug, &s.Status, &s.CreatedAt, &s.UpdatedAt)
	if err == pgx.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("set current space: %w", err)
	}
	return &s, nil
}

func (r *PostgresRepository) ListMembers(ctx context.Context, userID, spaceID string) ([]Member, error) {
	rows, err := r.pool.Query(ctx, `SELECT u.id, u.email, u.username, sm.role, sm.status, sm.created_at
		FROM space_members sm JOIN users u ON u.id = sm.user_id
		WHERE sm.space_id = $1 AND EXISTS (SELECT 1 FROM space_members access WHERE access.space_id = $1 AND access.user_id = $2 AND access.status = 'active')
		ORDER BY sm.created_at`, spaceID, userID)
	if err != nil {
		return nil, fmt.Errorf("list members: %w", err)
	}
	defer rows.Close()
	members := make([]Member, 0)
	for rows.Next() {
		var m Member
		if err := rows.Scan(&m.UserID, &m.Email, &m.Username, &m.Role, &m.Status, &m.CreatedAt); err != nil {
			return nil, fmt.Errorf("scan member: %w", err)
		}
		members = append(members, m)
	}
	return members, rows.Err()
}

func (r *PostgresRepository) Invite(ctx context.Context, userID, spaceID string, req InviteRequest) (bool, error) {
	tag, err := r.pool.Exec(ctx, `INSERT INTO space_invitations (space_id, email, role, invited_by)
		SELECT $1, LOWER($2), $3, $4 WHERE EXISTS (
			SELECT 1 FROM space_members WHERE space_id = $1 AND user_id = $4 AND role IN ('owner', 'admin') AND status = 'active'
		)`, spaceID, req.Email, req.Role, userID)
	return tag.RowsAffected() > 0, err
}

func (r *PostgresRepository) UpdateMemberRole(ctx context.Context, userID, spaceID, memberID, role string) (bool, error) {
	tag, err := r.pool.Exec(ctx, `UPDATE space_members target SET role = $4, updated_at = NOW()
		WHERE target.space_id = $2 AND target.user_id = $3 AND target.role <> 'owner'
		AND EXISTS (SELECT 1 FROM space_members actor WHERE actor.space_id = $2 AND actor.user_id = $1 AND actor.role IN ('owner', 'admin') AND actor.status = 'active')`, userID, spaceID, memberID, role)
	return tag.RowsAffected() > 0, err
}

func (r *PostgresRepository) RemoveMember(ctx context.Context, userID, spaceID, memberID string) (bool, error) {
	tag, err := r.pool.Exec(ctx, `UPDATE space_members target SET status = 'removed', updated_at = NOW()
		WHERE target.space_id = $2 AND target.user_id = $3 AND target.role <> 'owner'
		AND EXISTS (SELECT 1 FROM space_members actor WHERE actor.space_id = $2 AND actor.user_id = $1 AND actor.role IN ('owner', 'admin') AND actor.status = 'active')`, userID, spaceID, memberID)
	return tag.RowsAffected() > 0, err
}

func slugify(name string) string {
	name = strings.ToLower(strings.TrimSpace(name))
	var b strings.Builder
	for _, r := range name {
		if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') {
			b.WriteRune(r)
		} else if b.Len() > 0 {
			b.WriteByte('-')
		}
	}
	slug := strings.Trim(b.String(), "-")
	if slug == "" {
		return "workspace"
	}
	return slug
}
