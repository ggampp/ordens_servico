package repository

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/ggampp/ordens_servico/backend/internal/model"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// ErrNotFound is returned when a row does not exist.
var ErrNotFound = errors.New("not found")

// EmployeeRepository persists employees and their positions.
type EmployeeRepository struct {
	pool *pgxpool.Pool
}

// NewEmployeeRepository builds an EmployeeRepository.
func NewEmployeeRepository(pool *pgxpool.Pool) *EmployeeRepository {
	return &EmployeeRepository{pool: pool}
}

// Create inserts a new employee.
func (r *EmployeeRepository) Create(ctx context.Context, e *model.Employee) error {
	return r.pool.QueryRow(ctx, `
		INSERT INTO employees (code, name, email, phone, role, status)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, created_at, updated_at`,
		e.Code, e.Name, e.Email, e.Phone, e.Role, e.Status,
	).Scan(&e.ID, &e.CreatedAt, &e.UpdatedAt)
}

// NextCode returns the next six-digit employee code.
func (r *EmployeeRepository) NextCode(ctx context.Context) (string, error) {
	var next int64
	if err := r.pool.QueryRow(ctx, `
		SELECT COALESCE(MAX(code::integer), 0) + 1
		FROM employees
		WHERE code ~ '^[0-9]{6}$'`,
	).Scan(&next); err != nil {
		return "", err
	}
	if next > 999999 {
		return "", errors.New("employee code sequence exhausted")
	}
	return fmt.Sprintf("%06d", next), nil
}

// Update modifies an existing, non-deleted employee.
func (r *EmployeeRepository) Update(ctx context.Context, e *model.Employee) error {
	tag, err := r.pool.Exec(ctx, `
		UPDATE employees
		SET name = $2, email = $3, phone = $4, role = $5, status = $6, updated_at = now()
		WHERE id = $1 AND deleted = FALSE`,
		e.ID, e.Name, e.Email, e.Phone, e.Role, e.Status)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}

// SoftDelete performs a logical delete.
func (r *EmployeeRepository) SoftDelete(ctx context.Context, id int64) error {
	tag, err := r.pool.Exec(ctx,
		`UPDATE employees SET deleted = TRUE, status = 'inactive', updated_at = now()
		 WHERE id = $1 AND deleted = FALSE`, id)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}

// GetByID returns a single employee with its latest position hydrated.
func (r *EmployeeRepository) GetByID(ctx context.Context, id int64) (*model.Employee, error) {
	e := &model.Employee{}
	err := r.pool.QueryRow(ctx, `
		SELECT id, code, name, email, phone, role, status, created_at, updated_at
		FROM employees WHERE id = $1 AND deleted = FALSE`, id).
		Scan(&e.ID, &e.Code, &e.Name, &e.Email, &e.Phone, &e.Role, &e.Status, &e.CreatedAt, &e.UpdatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	pos, _ := r.LastPosition(ctx, id)
	e.LastPosition = pos
	return e, nil
}

// List returns a filtered, paginated set of employees.
func (r *EmployeeRepository) List(ctx context.Context, f model.EmployeeFilter) ([]model.Employee, int64, error) {
	where := []string{"deleted = FALSE"}
	args := []any{}
	idx := 1

	if f.Status != "" {
		where = append(where, "status = $"+itoa(idx))
		args = append(args, f.Status)
		idx++
	}
	if f.Search != "" {
		where = append(where, "(name ILIKE $"+itoa(idx)+" OR code ILIKE $"+itoa(idx)+")")
		args = append(args, "%"+f.Search+"%")
		idx++
	}
	cond := strings.Join(where, " AND ")

	var total int64
	if err := r.pool.QueryRow(ctx,
		"SELECT COUNT(*) FROM employees WHERE "+cond, args...).Scan(&total); err != nil {
		return nil, 0, err
	}

	args = append(args, f.PageSize, f.Offset())
	rows, err := r.pool.Query(ctx, `
		SELECT id, code, name, email, phone, role, status, created_at, updated_at
		FROM employees WHERE `+cond+`
		ORDER BY name ASC
		LIMIT $`+itoa(idx)+` OFFSET $`+itoa(idx+1), args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var out []model.Employee
	for rows.Next() {
		var e model.Employee
		if err := rows.Scan(&e.ID, &e.Code, &e.Name, &e.Email, &e.Phone,
			&e.Role, &e.Status, &e.CreatedAt, &e.UpdatedAt); err != nil {
			return nil, 0, err
		}
		out = append(out, e)
	}
	return out, total, rows.Err()
}

// AddPosition records a geolocation point for an employee.
func (r *EmployeeRepository) AddPosition(ctx context.Context, p *model.EmployeePosition) error {
	return r.pool.QueryRow(ctx, `
		INSERT INTO employee_positions (employee_id, latitude, longitude, geom, recorded_at)
		VALUES ($1, $2, $3, ST_SetSRID(ST_MakePoint($3, $2), 4326)::geography, $4)
		RETURNING id`,
		p.EmployeeID, p.Latitude, p.Longitude, p.RecordedAt).Scan(&p.ID)
}

// LastPosition returns the most recent position for an employee, if any.
func (r *EmployeeRepository) LastPosition(ctx context.Context, employeeID int64) (*model.EmployeePosition, error) {
	p := &model.EmployeePosition{}
	err := r.pool.QueryRow(ctx, `
		SELECT id, employee_id, latitude, longitude, recorded_at
		FROM employee_positions WHERE employee_id = $1
		ORDER BY recorded_at DESC LIMIT 1`, employeeID).
		Scan(&p.ID, &p.EmployeeID, &p.Latitude, &p.Longitude, &p.RecordedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return p, nil
}

// PositionHistory returns a paginated history of positions for an employee.
func (r *EmployeeRepository) PositionHistory(ctx context.Context, employeeID int64, p model.Pagination) ([]model.EmployeePosition, int64, error) {
	var total int64
	if err := r.pool.QueryRow(ctx,
		`SELECT COUNT(*) FROM employee_positions WHERE employee_id = $1`, employeeID).Scan(&total); err != nil {
		return nil, 0, err
	}
	rows, err := r.pool.Query(ctx, `
		SELECT id, employee_id, latitude, longitude, recorded_at
		FROM employee_positions WHERE employee_id = $1
		ORDER BY recorded_at DESC LIMIT $2 OFFSET $3`,
		employeeID, p.PageSize, p.Offset())
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()
	var out []model.EmployeePosition
	for rows.Next() {
		var ep model.EmployeePosition
		if err := rows.Scan(&ep.ID, &ep.EmployeeID, &ep.Latitude, &ep.Longitude, &ep.RecordedAt); err != nil {
			return nil, 0, err
		}
		out = append(out, ep)
	}
	return out, total, rows.Err()
}

// Exists reports whether an active employee exists.
func (r *EmployeeRepository) Exists(ctx context.Context, id int64) (bool, error) {
	var ok bool
	err := r.pool.QueryRow(ctx,
		`SELECT EXISTS(SELECT 1 FROM employees WHERE id = $1 AND deleted = FALSE)`, id).Scan(&ok)
	return ok, err
}
