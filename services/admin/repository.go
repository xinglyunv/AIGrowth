package admin

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"
)

const selectCols = `id, username, email, password_hash, role, status, created_at, updated_at`

type Repository interface {
	FindByID(ctx context.Context, id string) (*Admin, error)
	FindByEmail(ctx context.Context, email string) (*Admin, error)
	Create(ctx context.Context, req CreateAdminRequest) (*Admin, error)
	List(ctx context.Context, offset, limit int) ([]Admin, int, error)
	Update(ctx context.Context, id string, req UpdateAdminRequest) (*Admin, error)
	Delete(ctx context.Context, id string) error
	VerifyPassword(plainPassword, hashedPassword string) bool
	HashPassword(password string) (string, error)
}

type PostgresRepository struct {
	pool *pgxpool.Pool
}

func NewPostgresRepository(pool *pgxpool.Pool) Repository {
	return &PostgresRepository{pool: pool}
}

func (r *PostgresRepository) FindByID(ctx context.Context, id string) (*Admin, error) {
	a := &Admin{}
	err := r.pool.QueryRow(ctx,
		`SELECT `+selectCols+` FROM admins WHERE id = $1`, id,
	).Scan(&a.ID, &a.Username, &a.Email, &a.PasswordHash, &a.Role, &a.Status, &a.CreatedAt, &a.UpdatedAt)
	if err == pgx.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("find admin by id: %w", err)
	}
	return a, nil
}

func (r *PostgresRepository) FindByEmail(ctx context.Context, email string) (*Admin, error) {
	a := &Admin{}
	err := r.pool.QueryRow(ctx,
		`SELECT `+selectCols+` FROM admins WHERE email = $1`, email,
	).Scan(&a.ID, &a.Username, &a.Email, &a.PasswordHash, &a.Role, &a.Status, &a.CreatedAt, &a.UpdatedAt)
	if err == pgx.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("find admin by email: %w", err)
	}
	return a, nil
}

func (r *PostgresRepository) Create(ctx context.Context, req CreateAdminRequest) (*Admin, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("hash password: %w", err)
	}

	role := req.Role
	if role == "" {
		role = "admin"
	}

	a := &Admin{}
	err = r.pool.QueryRow(ctx,
		`INSERT INTO admins (username, email, password_hash, role)
		 VALUES ($1, $2, $3, $4)
		 RETURNING `+selectCols,
		req.Username, req.Email, string(hash), role,
	).Scan(&a.ID, &a.Username, &a.Email, &a.PasswordHash, &a.Role, &a.Status, &a.CreatedAt, &a.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("create admin: %w", err)
	}
	return a, nil
}

func (r *PostgresRepository) List(ctx context.Context, offset, limit int) ([]Admin, int, error) {
	var total int
	err := r.pool.QueryRow(ctx, `SELECT COUNT(*) FROM admins`).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("count admins: %w", err)
	}

	rows, err := r.pool.Query(ctx,
		`SELECT `+selectCols+` FROM admins ORDER BY created_at DESC LIMIT $1 OFFSET $2`, limit, offset,
	)
	if err != nil {
		return nil, 0, fmt.Errorf("list admins: %w", err)
	}
	defer rows.Close()

	var admins []Admin
	for rows.Next() {
		var a Admin
		err := rows.Scan(&a.ID, &a.Username, &a.Email, &a.PasswordHash, &a.Role, &a.Status, &a.CreatedAt, &a.UpdatedAt)
		if err != nil {
			return nil, 0, fmt.Errorf("scan admin: %w", err)
		}
		admins = append(admins, a)
	}
	if admins == nil {
		admins = []Admin{}
	}
	return admins, total, nil
}

func (r *PostgresRepository) Update(ctx context.Context, id string, req UpdateAdminRequest) (*Admin, error) {
	setClauses := []string{}
	args := []interface{}{}
	argIdx := 1

	if req.Username != "" {
		setClauses = append(setClauses, fmt.Sprintf("username = $%d", argIdx))
		args = append(args, req.Username)
		argIdx++
	}
	if req.Email != "" {
		setClauses = append(setClauses, fmt.Sprintf("email = $%d", argIdx))
		args = append(args, req.Email)
		argIdx++
	}
	if req.Password != "" {
		hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		if err != nil {
			return nil, fmt.Errorf("hash password: %w", err)
		}
		setClauses = append(setClauses, fmt.Sprintf("password_hash = $%d", argIdx))
		args = append(args, string(hash))
		argIdx++
	}
	if req.Role != "" {
		setClauses = append(setClauses, fmt.Sprintf("role = $%d", argIdx))
		args = append(args, req.Role)
		argIdx++
	}
	if req.Status != "" {
		setClauses = append(setClauses, fmt.Sprintf("status = $%d", argIdx))
		args = append(args, req.Status)
		argIdx++
	}

	if len(setClauses) == 0 {
		return r.FindByID(ctx, id)
	}

	setClauses = append(setClauses, fmt.Sprintf("updated_at = NOW()"))
	args = append(args, id)

	query := fmt.Sprintf(`UPDATE admins SET %s WHERE id = $%d RETURNING `+selectCols,
		joinStrings(setClauses, ", "), argIdx)

	a := &Admin{}
	err := r.pool.QueryRow(ctx, query, args...).Scan(
		&a.ID, &a.Username, &a.Email, &a.PasswordHash, &a.Role, &a.Status, &a.CreatedAt, &a.UpdatedAt,
	)
	if err == pgx.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("update admin: %w", err)
	}
	return a, nil
}

func (r *PostgresRepository) Delete(ctx context.Context, id string) error {
	_, err := r.pool.Exec(ctx, `DELETE FROM admins WHERE id = $1`, id)
	if err != nil {
		return fmt.Errorf("delete admin: %w", err)
	}
	return nil
}

func (r *PostgresRepository) VerifyPassword(plainPassword, hashedPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(plainPassword))
	return err == nil
}

func (r *PostgresRepository) HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
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
