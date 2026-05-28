package repository

import (
	"context"

	"github.com/ggampp/ordens_servico/backend/internal/model"
	"github.com/jackc/pgx/v5/pgxpool"
)

// DashboardRepository computes aggregate indicators.
type DashboardRepository struct {
	pool *pgxpool.Pool
}

// NewDashboardRepository builds a DashboardRepository.
func NewDashboardRepository(pool *pgxpool.Pool) *DashboardRepository {
	return &DashboardRepository{pool: pool}
}

// Summary returns the dashboard indicators in a few aggregate queries.
func (r *DashboardRepository) Summary(ctx context.Context) (*model.DashboardSummary, error) {
	s := &model.DashboardSummary{OrdersByPriority: map[string]int64{}}

	// Counts per status.
	rows, err := r.pool.Query(ctx,
		`SELECT status, COUNT(*) FROM service_orders WHERE deleted = FALSE GROUP BY status`)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var status string
		var n int64
		if err := rows.Scan(&status, &n); err != nil {
			rows.Close()
			return nil, err
		}
		switch status {
		case model.StatusOpen:
			s.OpenOrders = n
		case model.StatusAssigned:
			s.AssignedOrders = n
		case model.StatusInProgress:
			s.InProgressOrders = n
		case model.StatusCompleted:
			s.CompletedOrders = n
		case model.StatusCancelled:
			s.CancelledOrders = n
		}
	}
	rows.Close()
	if err := rows.Err(); err != nil {
		return nil, err
	}

	// Counts per priority.
	prows, err := r.pool.Query(ctx,
		`SELECT priority, COUNT(*) FROM service_orders WHERE deleted = FALSE GROUP BY priority`)
	if err != nil {
		return nil, err
	}
	for prows.Next() {
		var priority string
		var n int64
		if err := prows.Scan(&priority, &n); err != nil {
			prows.Close()
			return nil, err
		}
		s.OrdersByPriority[priority] = n
	}
	prows.Close()
	if err := prows.Err(); err != nil {
		return nil, err
	}

	// Active employees.
	if err := r.pool.QueryRow(ctx,
		`SELECT COUNT(*) FROM employees WHERE deleted = FALSE AND status = 'active'`).
		Scan(&s.ActiveEmployees); err != nil {
		return nil, err
	}

	// Orders grouped by responsible employee.
	erows, err := r.pool.Query(ctx, `
		SELECT o.employee_id, COALESCE(e.name, 'Não atribuído') AS name, COUNT(*) AS total
		FROM service_orders o
		LEFT JOIN employees e ON e.id = o.employee_id
		WHERE o.deleted = FALSE
		GROUP BY o.employee_id, e.name
		ORDER BY total DESC`)
	if err != nil {
		return nil, err
	}
	defer erows.Close()
	for erows.Next() {
		var item model.OrdersByEmployee
		if err := erows.Scan(&item.EmployeeID, &item.EmployeeName, &item.Count); err != nil {
			return nil, err
		}
		s.OrdersByEmployee = append(s.OrdersByEmployee, item)
	}
	return s, erows.Err()
}
