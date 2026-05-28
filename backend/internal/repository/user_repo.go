package repository

import (
	"context"
	"errors"

	"github.com/ggampp/ordens_servico/backend/internal/model"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// UserRepository persists authentication users.
type UserRepository struct {
	pool *pgxpool.Pool
}

// NewUserRepository builds a UserRepository.
func NewUserRepository(pool *pgxpool.Pool) *UserRepository {
	return &UserRepository{pool: pool}
}

// Create inserts a user.
func (r *UserRepository) Create(ctx context.Context, u *model.User) error {
	return r.pool.QueryRow(ctx, `
		INSERT INTO users (name, email, password_hash, role, employee_id, active)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, created_at, updated_at`,
		u.Name, u.Email, u.PasswordHash, u.Role, u.EmployeeID, u.Active).
		Scan(&u.ID, &u.CreatedAt, &u.UpdatedAt)
}

// GetByEmail returns the user matching an email.
func (r *UserRepository) GetByEmail(ctx context.Context, email string) (*model.User, error) {
	u := &model.User{}
	err := r.pool.QueryRow(ctx, `
		SELECT id, name, email, password_hash, role, employee_id, active, created_at, updated_at
		FROM users WHERE email = $1`, email).
		Scan(&u.ID, &u.Name, &u.Email, &u.PasswordHash, &u.Role, &u.EmployeeID,
			&u.Active, &u.CreatedAt, &u.UpdatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return u, nil
}

// CountAdmins returns the number of active admin users, used for seeding.
func (r *UserRepository) CountAdmins(ctx context.Context) (int64, error) {
	var n int64
	err := r.pool.QueryRow(ctx,
		`SELECT COUNT(*) FROM users WHERE role = 'admin'`).Scan(&n)
	return n, err
}
