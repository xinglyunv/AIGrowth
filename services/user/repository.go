package user

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Repository defines the user data access interface.
type Repository interface {
	Create(ctx context.Context, req CreateUserRequest) (*User, error)
	FindByEmail(ctx context.Context, email string) (*User, error)
	FindByPhone(ctx context.Context, phone string) (*User, error)
	FindByID(ctx context.Context, id string) (*User, error)
	Update(ctx context.Context, id string, req UpdateUserRequest) (*User, error)
	UpdatePassword(ctx context.Context, id string, hash string) error
	UpdateLastLogin(ctx context.Context, id string) error
	List(ctx context.Context, offset, limit int) ([]User, int, error)
	SaveVerificationCode(ctx context.Context, email, code, purpose string) error
	VerifyCode(ctx context.Context, email, code, purpose string) (bool, error)
	MarkCodeUsed(ctx context.Context, email, code, purpose string) error
	AddCredits(ctx context.Context, id string, amount int) (*User, error)
	DeductCredits(ctx context.Context, id string, amount int) (*User, error)
}

// PostgresRepository implements Repository using pgxpool.
type PostgresRepository struct {
	pool *pgxpool.Pool
}

// NewPostgresRepository creates a new PostgresRepository.
func NewPostgresRepository(pool *pgxpool.Pool) Repository {
	return &PostgresRepository{pool: pool}
}

// Create inserts a new user.
func (r *PostgresRepository) Create(ctx context.Context, req CreateUserRequest) (*User, error) {
	u := &User{}
	err := r.pool.QueryRow(ctx,
		`INSERT INTO users (email, phone, password_hash, username, company_name)
		 VALUES ($1, $2, $3, $4, $5)
		 RETURNING id, email, phone, password_hash, username, company_name, avatar_url,
		           role, email_verified, status, COALESCE(credits, 0) AS credits, last_login_at, created_at, updated_at`,
		req.Email, req.Phone, req.Password, req.Username, req.CompanyName,
	).Scan(
		&u.ID, &u.Email, &u.Phone, &u.PasswordHash, &u.Username, &u.CompanyName,
		&u.AvatarURL, &u.Role, &u.EmailVerified, &u.Status, &u.Credits, &u.LastLoginAt,
		&u.CreatedAt, &u.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("create user: %w", err)
	}
	return u, nil
}

// FindByEmail looks up a user by email.
func (r *PostgresRepository) FindByEmail(ctx context.Context, email string) (*User, error) {
	u := &User{}
	err := r.pool.QueryRow(ctx,
		`SELECT id, email, phone, password_hash, username, company_name, avatar_url,
		        role, email_verified, status, COALESCE(credits, 0) AS credits, last_login_at, created_at, updated_at
		 FROM users WHERE email = $1`, email,
	).Scan(
		&u.ID, &u.Email, &u.Phone, &u.PasswordHash, &u.Username, &u.CompanyName,
		&u.AvatarURL, &u.Role, &u.EmailVerified, &u.Status, &u.Credits, &u.LastLoginAt,
		&u.CreatedAt, &u.UpdatedAt,
	)
	if err == pgx.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("find user by email: %w", err)
	}
	return u, nil
}

// FindByPhone looks up a user by phone number.
func (r *PostgresRepository) FindByPhone(ctx context.Context, phone string) (*User, error) {
	u := &User{}
	err := r.pool.QueryRow(ctx,
		`SELECT id, email, phone, password_hash, username, company_name, avatar_url,
		        role, email_verified, status, COALESCE(credits, 0) AS credits, last_login_at, created_at, updated_at
		 FROM users WHERE phone = $1`, phone,
	).Scan(
		&u.ID, &u.Email, &u.Phone, &u.PasswordHash, &u.Username, &u.CompanyName,
		&u.AvatarURL, &u.Role, &u.EmailVerified, &u.Status, &u.Credits, &u.LastLoginAt,
		&u.CreatedAt, &u.UpdatedAt,
	)
	if err == pgx.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("find user by phone: %w", err)
	}
	return u, nil
}

// FindByID looks up a user by ID.
func (r *PostgresRepository) FindByID(ctx context.Context, id string) (*User, error) {
	u := &User{}
	err := r.pool.QueryRow(ctx,
		`SELECT id, email, phone, password_hash, username, company_name, avatar_url,
		        role, email_verified, status, COALESCE(credits, 0) AS credits, last_login_at, created_at, updated_at
		 FROM users WHERE id = $1`, id,
	).Scan(
		&u.ID, &u.Email, &u.Phone, &u.PasswordHash, &u.Username, &u.CompanyName,
		&u.AvatarURL, &u.Role, &u.EmailVerified, &u.Status, &u.Credits, &u.LastLoginAt,
		&u.CreatedAt, &u.UpdatedAt,
	)
	if err == pgx.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("find user by id: %w", err)
	}
	return u, nil
}

// Update modifies user profile fields.
func (r *PostgresRepository) Update(ctx context.Context, id string, req UpdateUserRequest) (*User, error) {
	u := &User{}
	err := r.pool.QueryRow(ctx,
		`UPDATE users
		 SET username = COALESCE($2, username),
		     company_name = COALESCE($3, company_name),
		     avatar_url = COALESCE($4, avatar_url)
		 WHERE id = $1
		 RETURNING id, email, phone, password_hash, username, company_name, avatar_url,
		           role, email_verified, status, COALESCE(credits, 0) AS credits, last_login_at, created_at, updated_at`,
		id, req.Username, req.CompanyName, req.AvatarURL,
	).Scan(
		&u.ID, &u.Email, &u.Phone, &u.PasswordHash, &u.Username, &u.CompanyName,
		&u.AvatarURL, &u.Role, &u.EmailVerified, &u.Status, &u.Credits, &u.LastLoginAt,
		&u.CreatedAt, &u.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("update user: %w", err)
	}
	return u, nil
}

// UpdatePassword changes the user's password hash.
func (r *PostgresRepository) UpdatePassword(ctx context.Context, id string, hash string) error {
	_, err := r.pool.Exec(ctx, `UPDATE users SET password_hash = $2 WHERE id = $1`, id, hash)
	if err != nil {
		return fmt.Errorf("update password: %w", err)
	}
	return nil
}

// UpdateLastLogin sets the user's last login timestamp.
func (r *PostgresRepository) UpdateLastLogin(ctx context.Context, id string) error {
	now := time.Now()
	_, err := r.pool.Exec(ctx, `UPDATE users SET last_login_at = $2 WHERE id = $1`, id, now)
	if err != nil {
		return fmt.Errorf("update last login: %w", err)
	}
	return nil
}

// List returns paginated users with a total count.
func (r *PostgresRepository) List(ctx context.Context, offset, limit int) ([]User, int, error) {
	var total int
	err := r.pool.QueryRow(ctx, `SELECT COUNT(*) FROM users`).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("count users: %w", err)
	}

	rows, err := r.pool.Query(ctx,
		`SELECT id, email, phone, password_hash, username, company_name, avatar_url,
		        role, email_verified, status, COALESCE(credits, 0) AS credits, last_login_at, created_at, updated_at
		 FROM users ORDER BY created_at DESC LIMIT $1 OFFSET $2`, limit, offset,
	)
	if err != nil {
		return nil, 0, fmt.Errorf("list users: %w", err)
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		var u User
		err := rows.Scan(
			&u.ID, &u.Email, &u.Phone, &u.PasswordHash, &u.Username, &u.CompanyName,
			&u.AvatarURL, &u.Role, &u.EmailVerified, &u.Status, &u.Credits, &u.LastLoginAt,
			&u.CreatedAt, &u.UpdatedAt,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("scan user: %w", err)
		}
		users = append(users, u)
	}
	if users == nil {
		users = []User{}
	}
	return users, total, nil
}

// AddCredits adds credits to a user's balance and returns the updated user.
func (r *PostgresRepository) AddCredits(ctx context.Context, id string, amount int) (*User, error) {
	u := &User{}
	err := r.pool.QueryRow(ctx,
		`UPDATE users SET credits = credits + $2 WHERE id = $1
		 RETURNING id, email, phone, password_hash, username, company_name, avatar_url,
		           role, email_verified, status, COALESCE(credits, 0) AS credits, last_login_at, created_at, updated_at`,
		id, amount,
	).Scan(
		&u.ID, &u.Email, &u.Phone, &u.PasswordHash, &u.Username, &u.CompanyName,
		&u.AvatarURL, &u.Role, &u.EmailVerified, &u.Status, &u.Credits, &u.LastLoginAt,
		&u.CreatedAt, &u.UpdatedAt,
	)
	if err == pgx.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("add credits: %w", err)
	}
	return u, nil
}

// DeductCredits deducts credits from a user's balance and returns the updated user.
func (r *PostgresRepository) DeductCredits(ctx context.Context, id string, amount int) (*User, error) {
	u := &User{}
	err := r.pool.QueryRow(ctx,
		`UPDATE users SET credits = credits - $2, updated_at = NOW() WHERE id = $1 AND credits >= $2
		 RETURNING id, email, phone, password_hash, username, company_name, avatar_url,
		           role, email_verified, status, COALESCE(credits, 0) AS credits, last_login_at, created_at, updated_at`,
		id, amount,
	).Scan(
		&u.ID, &u.Email, &u.Phone, &u.PasswordHash, &u.Username, &u.CompanyName,
		&u.AvatarURL, &u.Role, &u.EmailVerified, &u.Status, &u.Credits, &u.LastLoginAt,
		&u.CreatedAt, &u.UpdatedAt,
	)
	if err == pgx.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("deduct credits: %w", err)
	}
	return u, nil
}

// SaveVerificationCode inserts a verification code record.
func (r *PostgresRepository) SaveVerificationCode(ctx context.Context, email, code, purpose string) error {
	expiresAt := time.Now().Add(10 * time.Minute)
	_, err := r.pool.Exec(ctx,
		`INSERT INTO verification_codes (email, code, purpose, expires_at) VALUES ($1, $2, $3, $4)`,
		email, code, purpose, expiresAt,
	)
	if err != nil {
		return fmt.Errorf("save verification code: %w", err)
	}
	return nil
}

// VerifyCode checks if the code is valid (unused and not expired).
func (r *PostgresRepository) VerifyCode(ctx context.Context, email, code, purpose string) (bool, error) {
	var valid bool
	err := r.pool.QueryRow(ctx,
		`SELECT EXISTS(
			SELECT 1 FROM verification_codes
			WHERE email = $1 AND code = $2 AND purpose = $3
			  AND used = false AND expires_at > NOW()
		)`, email, code, purpose,
	).Scan(&valid)
	if err != nil {
		return false, fmt.Errorf("verify code: %w", err)
	}
	return valid, nil
}

// MarkCodeUsed marks a verification code as used.
func (r *PostgresRepository) MarkCodeUsed(ctx context.Context, email, code, purpose string) error {
	_, err := r.pool.Exec(ctx,
		`UPDATE verification_codes SET used = true
		 WHERE email = $1 AND code = $2 AND purpose = $3 AND used = false`,
		email, code, purpose,
	)
	if err != nil {
		return fmt.Errorf("mark code used: %w", err)
	}
	return nil
}
