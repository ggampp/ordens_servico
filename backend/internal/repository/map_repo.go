package repository

import (
	"context"

	"github.com/ggampp/ordens_servico/backend/internal/model"
	"github.com/jackc/pgx/v5/pgxpool"
)

// MapRepository serves geolocation views.
type MapRepository struct {
	pool *pgxpool.Pool
}

// NewMapRepository builds a MapRepository.
func NewMapRepository(pool *pgxpool.Pool) *MapRepository {
	return &MapRepository{pool: pool}
}

// EmployeesWithPosition returns active employees that have at least one
// recorded position, hydrated with their latest location.
func (r *MapRepository) EmployeesWithPosition(ctx context.Context) ([]model.Employee, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT DISTINCT ON (e.id)
		       e.id, e.code, e.name, e.email, e.phone, e.role, e.status,
		       e.created_at, e.updated_at,
		       p.id, p.latitude, p.longitude, p.recorded_at
		FROM employees e
		JOIN employee_positions p ON p.employee_id = e.id
		WHERE e.deleted = FALSE AND e.status = 'active'
		ORDER BY e.id, p.recorded_at DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []model.Employee
	for rows.Next() {
		var e model.Employee
		var pos model.EmployeePosition
		if err := rows.Scan(&e.ID, &e.Code, &e.Name, &e.Email, &e.Phone, &e.Role,
			&e.Status, &e.CreatedAt, &e.UpdatedAt,
			&pos.ID, &pos.Latitude, &pos.Longitude, &pos.RecordedAt); err != nil {
			return nil, err
		}
		pos.EmployeeID = e.ID
		e.LastPosition = &pos
		out = append(out, e)
	}
	return out, rows.Err()
}
