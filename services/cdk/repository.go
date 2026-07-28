package cdk

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

const cdkCols = `id, code, credits, max_uses, used_count, is_active, expires_at, created_at, updated_at`

type Repository interface {
	List(ctx context.Context, offset, limit int) ([]CDKCode, int, error)
	GetByID(ctx context.Context, id string) (*CDKCode, error)
	GetByCode(ctx context.Context, code string) (*CDKCode, error)
	Create(ctx context.Context, req CreateCDKRequest) (*CDKCode, error)
	Update(ctx context.Context, id string, req UpdateCDKRequest) (*CDKCode, error)
	Redeem(ctx context.Context, code, userID string) (*RedeemResult, error)
	GetUsages(ctx context.Context, cdkID string) ([]CDKUsage, error)
}

type PostgresRepository struct {
	pool *pgxpool.Pool
}

func NewPostgresRepository(pool *pgxpool.Pool) Repository {
	return &PostgresRepository{pool: pool}
}

func (r *PostgresRepository) List(ctx context.Context, offset, limit int) ([]CDKCode, int, error) {
	var total int
	err := r.pool.QueryRow(ctx, `SELECT COUNT(*) FROM cdk_codes`).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("count cdks: %w", err)
	}

	rows, err := r.pool.Query(ctx,
		`SELECT `+cdkCols+` FROM cdk_codes ORDER BY created_at DESC LIMIT $1 OFFSET $2`,
		limit, offset,
	)
	if err != nil {
		return nil, 0, fmt.Errorf("list cdks: %w", err)
	}
	defer rows.Close()

	var codes []CDKCode
	for rows.Next() {
		var c CDKCode
		err := rows.Scan(&c.ID, &c.Code, &c.Credits, &c.MaxUses, &c.UsedCount, &c.IsActive, &c.ExpiresAt, &c.CreatedAt, &c.UpdatedAt)
		if err != nil {
			return nil, 0, fmt.Errorf("scan cdk: %w", err)
		}
		codes = append(codes, c)
	}
	if codes == nil {
		codes = []CDKCode{}
	}
	return codes, total, nil
}

func (r *PostgresRepository) GetByID(ctx context.Context, id string) (*CDKCode, error) {
	c := &CDKCode{}
	err := r.pool.QueryRow(ctx,
		`SELECT `+cdkCols+` FROM cdk_codes WHERE id = $1`, id,
	).Scan(&c.ID, &c.Code, &c.Credits, &c.MaxUses, &c.UsedCount, &c.IsActive, &c.ExpiresAt, &c.CreatedAt, &c.UpdatedAt)
	if err == pgx.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("get cdk by id: %w", err)
	}
	return c, nil
}

func (r *PostgresRepository) GetByCode(ctx context.Context, code string) (*CDKCode, error) {
	c := &CDKCode{}
	err := r.pool.QueryRow(ctx,
		`SELECT `+cdkCols+` FROM cdk_codes WHERE code = $1`, code,
	).Scan(&c.ID, &c.Code, &c.Credits, &c.MaxUses, &c.UsedCount, &c.IsActive, &c.ExpiresAt, &c.CreatedAt, &c.UpdatedAt)
	if err == pgx.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("get cdk by code: %w", err)
	}
	return c, nil
}

func (r *PostgresRepository) Create(ctx context.Context, req CreateCDKRequest) (*CDKCode, error) {
	if strings.TrimSpace(req.Code) == "" {
		buf := make([]byte, 8)
		if _, err := rand.Read(buf); err != nil {
			return nil, fmt.Errorf("generate cdk code: %w", err)
		}
		req.Code = "AIGE-" + strings.ToUpper(hex.EncodeToString(buf))
	} else {
		req.Code = strings.ToUpper(strings.TrimSpace(req.Code))
	}
	if req.Credits <= 0 {
		return nil, fmt.Errorf("CDK 额度必须大于 0")
	}
	if req.MaxUses < 0 {
		return nil, fmt.Errorf("CDK 最大使用次数不能小于 0")
	}
	c := &CDKCode{}
	err := r.pool.QueryRow(ctx,
		`INSERT INTO cdk_codes (code, credits, max_uses, expires_at)
		 VALUES ($1, $2, $3, $4)
		 RETURNING `+cdkCols,
		req.Code, req.Credits, req.MaxUses, req.ExpiresAt,
	).Scan(&c.ID, &c.Code, &c.Credits, &c.MaxUses, &c.UsedCount, &c.IsActive, &c.ExpiresAt, &c.CreatedAt, &c.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("create cdk: %w", err)
	}
	return c, nil
}

func (r *PostgresRepository) Update(ctx context.Context, id string, req UpdateCDKRequest) (*CDKCode, error) {
	setClauses := []string{}
	args := []interface{}{}
	argIdx := 1

	if req.IsActive != nil {
		setClauses = append(setClauses, fmt.Sprintf("is_active = $%d", argIdx))
		args = append(args, *req.IsActive)
		argIdx++
	}
	if req.MaxUses != nil {
		setClauses = append(setClauses, fmt.Sprintf("max_uses = $%d", argIdx))
		args = append(args, *req.MaxUses)
		argIdx++
	}
	if req.ExpiresAt != nil {
		setClauses = append(setClauses, fmt.Sprintf("expires_at = $%d", argIdx))
		args = append(args, *req.ExpiresAt)
		argIdx++
	}

	if len(setClauses) == 0 {
		return r.GetByID(ctx, id)
	}

	setClauses = append(setClauses, "updated_at = NOW()")
	args = append(args, id)

	query := fmt.Sprintf(`UPDATE cdk_codes SET %s WHERE id = $%d RETURNING `+cdkCols,
		joinStrings(setClauses, ", "), argIdx)

	c := &CDKCode{}
	err := r.pool.QueryRow(ctx, query, args...).Scan(
		&c.ID, &c.Code, &c.Credits, &c.MaxUses, &c.UsedCount, &c.IsActive, &c.ExpiresAt, &c.CreatedAt, &c.UpdatedAt,
	)
	if err == pgx.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("update cdk: %w", err)
	}
	return c, nil
}

func (r *PostgresRepository) Redeem(ctx context.Context, code, userID string) (*RedeemResult, error) {
	cdk, err := r.GetByCode(ctx, code)
	if err != nil {
		return nil, fmt.Errorf("get cdk: %w", err)
	}
	if cdk == nil {
		return &RedeemResult{Credits: 0, Message: "CDK code not found"}, nil
	}
	if !cdk.IsActive {
		return &RedeemResult{Credits: 0, Message: "CDK code is inactive"}, nil
	}
	if cdk.ExpiresAt != nil && cdk.ExpiresAt.Before(time.Now()) {
		return &RedeemResult{Credits: 0, Message: "CDK code has expired"}, nil
	}
	if cdk.UsedCount >= cdk.MaxUses {
		return &RedeemResult{Credits: 0, Message: "CDK code has been fully used"}, nil
	}

	var usedCount int
	err = r.pool.QueryRow(ctx,
		`SELECT COUNT(*) FROM cdk_usages WHERE cdk_id = $1 AND user_id = $2`,
		cdk.ID, userID,
	).Scan(&usedCount)
	if err != nil {
		return nil, fmt.Errorf("check cdk usage: %w", err)
	}
	if usedCount > 0 {
		return &RedeemResult{Credits: 0, Message: "You have already used this CDK code"}, nil
	}

	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("begin tx: %w", err)
	}
	defer tx.Rollback(ctx)

	_, err = tx.Exec(ctx,
		`UPDATE cdk_codes SET used_count = used_count + 1, updated_at = NOW() WHERE id = $1`, cdk.ID)
	if err != nil {
		return nil, fmt.Errorf("increment cdk count: %w", err)
	}

	_, err = tx.Exec(ctx,
		`INSERT INTO cdk_usages (cdk_id, user_id) VALUES ($1, $2)`, cdk.ID, userID)
	if err != nil {
		return nil, fmt.Errorf("record cdk usage: %w", err)
	}

	_, err = tx.Exec(ctx,
		`UPDATE users SET credits = credits + $1 WHERE id = $2`, cdk.Credits, userID)
	if err != nil {
		return nil, fmt.Errorf("add credits: %w", err)
	}

	err = tx.Commit(ctx)
	if err != nil {
		return nil, fmt.Errorf("commit tx: %w", err)
	}

	return &RedeemResult{Credits: cdk.Credits, Message: "Redeemed successfully"}, nil
}

func (r *PostgresRepository) GetUsages(ctx context.Context, cdkID string) ([]CDKUsage, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, cdk_id, user_id, used_at FROM cdk_usages WHERE cdk_id = $1 ORDER BY used_at DESC`, cdkID)
	if err != nil {
		return nil, fmt.Errorf("get cdk usages: %w", err)
	}
	defer rows.Close()

	var usages []CDKUsage
	for rows.Next() {
		var u CDKUsage
		err := rows.Scan(&u.ID, &u.CDKID, &u.UserID, &u.UsedAt)
		if err != nil {
			return nil, fmt.Errorf("scan cdk usage: %w", err)
		}
		usages = append(usages, u)
	}
	if usages == nil {
		usages = []CDKUsage{}
	}
	return usages, nil
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
