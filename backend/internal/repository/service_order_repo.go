package repository

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/ggampp/ordens_servico/backend/internal/model"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// ServiceOrderRepository persists service orders and their status history.
type ServiceOrderRepository struct {
	pool *pgxpool.Pool
}

// NewServiceOrderRepository builds a ServiceOrderRepository.
func NewServiceOrderRepository(pool *pgxpool.Pool) *ServiceOrderRepository {
	return &ServiceOrderRepository{pool: pool}
}

const orderSelect = `
	SELECT o.id, o.number, o.title, o.description, o.priority, o.status,
	       o.employee_id, e.name, o.address, o.latitude, o.longitude,
	       o.opened_at, o.due_at, o.completed_at, o.notes, o.created_at, o.updated_at
	FROM service_orders o
	LEFT JOIN employees e ON e.id = o.employee_id`

func scanOrder(row pgx.Row, o *model.ServiceOrder) error {
	return row.Scan(&o.ID, &o.Number, &o.Title, &o.Description, &o.Priority, &o.Status,
		&o.EmployeeID, &o.EmployeeName, &o.Address, &o.Latitude, &o.Longitude,
		&o.OpenedAt, &o.DueAt, &o.CompletedAt, &o.Notes, &o.CreatedAt, &o.UpdatedAt)
}

// Create inserts a new order and its initial history row in one transaction.
func (r *ServiceOrderRepository) Create(ctx context.Context, o *model.ServiceOrder, changedBy *int64) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	err = tx.QueryRow(ctx, `
		INSERT INTO service_orders
			(number, title, description, priority, status, employee_id, address,
			 latitude, longitude, geom, due_at, notes)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8::double precision,$9::double precision,
			CASE WHEN $8::double precision IS NULL OR $9::double precision IS NULL THEN NULL
			     ELSE ST_SetSRID(ST_MakePoint($9::double precision, $8::double precision), 4326)::geography END,
			$10,$11)
		RETURNING id, opened_at, created_at, updated_at`,
		o.Number, o.Title, o.Description, o.Priority, o.Status, o.EmployeeID,
		o.Address, o.Latitude, o.Longitude, o.DueAt, o.Notes).
		Scan(&o.ID, &o.OpenedAt, &o.CreatedAt, &o.UpdatedAt)
	if err != nil {
		return err
	}

	if _, err := tx.Exec(ctx, `
		INSERT INTO service_order_history (service_order_id, old_status, new_status, changed_by, note)
		VALUES ($1, NULL, $2, $3, 'created')`,
		o.ID, o.Status, changedBy); err != nil {
		return err
	}
	return tx.Commit(ctx)
}

// Update edits mutable fields of a non-deleted order.
func (r *ServiceOrderRepository) Update(ctx context.Context, o *model.ServiceOrder) error {
	tag, err := r.pool.Exec(ctx, `
		UPDATE service_orders SET
			title = $2, description = $3, priority = $4, address = $5,
			latitude = $6::double precision, longitude = $7::double precision,
			geom = CASE WHEN $6::double precision IS NULL OR $7::double precision IS NULL THEN NULL
			            ELSE ST_SetSRID(ST_MakePoint($7::double precision, $6::double precision), 4326)::geography END,
			due_at = $8, notes = $9, updated_at = now()
		WHERE id = $1 AND deleted = FALSE`,
		o.ID, o.Title, o.Description, o.Priority, o.Address,
		o.Latitude, o.Longitude, o.DueAt, o.Notes)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}

// SoftDelete logically removes an order.
func (r *ServiceOrderRepository) SoftDelete(ctx context.Context, id int64) error {
	tag, err := r.pool.Exec(ctx,
		`UPDATE service_orders SET deleted = TRUE, updated_at = now()
		 WHERE id = $1 AND deleted = FALSE`, id)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}

// GetByID returns a single order.
func (r *ServiceOrderRepository) GetByID(ctx context.Context, id int64) (*model.ServiceOrder, error) {
	o := &model.ServiceOrder{}
	err := scanOrder(r.pool.QueryRow(ctx, orderSelect+" WHERE o.id = $1 AND o.deleted = FALSE", id), o)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return o, nil
}

// ChangeStatus updates status and records a history entry transactionally.
// completed_at is set/cleared automatically based on the target status.
func (r *ServiceOrderRepository) ChangeStatus(ctx context.Context, id int64, oldStatus, newStatus string, changedBy *int64, note string) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	// completed_at is derived from the target status in Go to avoid reusing the
	// status parameter in two type contexts within the same statement.
	var completedAt *time.Time
	if newStatus == model.StatusCompleted {
		now := time.Now()
		completedAt = &now
	}
	tag, err := tx.Exec(ctx, `
		UPDATE service_orders SET
			status = $2,
			completed_at = $3,
			updated_at = now()
		WHERE id = $1 AND deleted = FALSE`, id, newStatus, completedAt)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return ErrNotFound
	}

	var notePtr *string
	if note != "" {
		notePtr = &note
	}
	if _, err := tx.Exec(ctx, `
		INSERT INTO service_order_history (service_order_id, old_status, new_status, changed_by, note)
		VALUES ($1, $2, $3, $4, $5)`,
		id, oldStatus, newStatus, changedBy, notePtr); err != nil {
		return err
	}
	return tx.Commit(ctx)
}

// Assign sets the responsible employee, moving an open order to "assigned".
func (r *ServiceOrderRepository) Assign(ctx context.Context, id, employeeID int64, currentStatus string, changedBy *int64) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	newStatus := currentStatus
	if currentStatus == model.StatusOpen {
		newStatus = model.StatusAssigned
	}

	tag, err := tx.Exec(ctx, `
		UPDATE service_orders SET employee_id = $2, status = $3, updated_at = now()
		WHERE id = $1 AND deleted = FALSE`, id, employeeID, newStatus)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return ErrNotFound
	}

	if newStatus != currentStatus {
		if _, err := tx.Exec(ctx, `
			INSERT INTO service_order_history (service_order_id, old_status, new_status, changed_by, note)
			VALUES ($1, $2, $3, $4, 'assigned')`,
			id, currentStatus, newStatus, changedBy); err != nil {
			return err
		}
	}
	return tx.Commit(ctx)
}

// History returns the status change history for an order.
func (r *ServiceOrderRepository) History(ctx context.Context, orderID int64) ([]model.ServiceOrderHistory, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT h.id, h.service_order_id, h.old_status, h.new_status,
		       h.changed_by, u.name, h.note, h.changed_at
		FROM service_order_history h
		LEFT JOIN users u ON u.id = h.changed_by
		WHERE h.service_order_id = $1
		ORDER BY h.changed_at ASC`, orderID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []model.ServiceOrderHistory
	for rows.Next() {
		var h model.ServiceOrderHistory
		if err := rows.Scan(&h.ID, &h.ServiceOrderID, &h.OldStatus, &h.NewStatus,
			&h.ChangedBy, &h.ChangedByName, &h.Note, &h.ChangedAt); err != nil {
			return nil, err
		}
		out = append(out, h)
	}
	return out, rows.Err()
}

// buildFilter constructs the shared WHERE clause for list/map queries.
func buildFilter(f model.ServiceOrderFilter, startIdx int) (string, []any, int) {
	where := []string{"o.deleted = FALSE"}
	args := []any{}
	idx := startIdx

	add := func(clause string, val any) {
		where = append(where, fmt.Sprintf(clause, idx))
		args = append(args, val)
		idx++
	}
	if f.Status != "" {
		add("o.status = $%d", f.Status)
	}
	if f.Priority != "" {
		add("o.priority = $%d", f.Priority)
	}
	if f.EmployeeID != nil {
		add("o.employee_id = $%d", *f.EmployeeID)
	}
	if f.DateFrom != nil {
		add("o.opened_at >= $%d", *f.DateFrom)
	}
	if f.DateTo != nil {
		add("o.opened_at <= $%d", *f.DateTo)
	}
	if f.MinLat != nil && f.MinLng != nil && f.MaxLat != nil && f.MaxLng != nil {
		where = append(where, fmt.Sprintf(
			"o.geom && ST_MakeEnvelope($%d, $%d, $%d, $%d, 4326)::geography",
			idx, idx+1, idx+2, idx+3))
		args = append(args, *f.MinLng, *f.MinLat, *f.MaxLng, *f.MaxLat)
		idx += 4
	}
	return strings.Join(where, " AND "), args, idx
}

// List returns filtered, paginated orders.
func (r *ServiceOrderRepository) List(ctx context.Context, f model.ServiceOrderFilter) ([]model.ServiceOrder, int64, error) {
	cond, args, idx := buildFilter(f, 1)

	var total int64
	if err := r.pool.QueryRow(ctx,
		"SELECT COUNT(*) FROM service_orders o WHERE "+cond, args...).Scan(&total); err != nil {
		return nil, 0, err
	}

	args = append(args, f.PageSize, f.Offset())
	query := orderSelect + " WHERE " + cond +
		fmt.Sprintf(" ORDER BY o.opened_at DESC LIMIT $%d OFFSET $%d", idx, idx+1)
	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()
	var out []model.ServiceOrder
	for rows.Next() {
		var o model.ServiceOrder
		if err := scanOrder(rows, &o); err != nil {
			return nil, 0, err
		}
		out = append(out, o)
	}
	return out, total, rows.Err()
}

// ForMap returns orders that have coordinates, applying the same filters.
func (r *ServiceOrderRepository) ForMap(ctx context.Context, f model.ServiceOrderFilter) ([]model.ServiceOrder, error) {
	cond, args, _ := buildFilter(f, 1)
	cond += " AND o.latitude IS NOT NULL AND o.longitude IS NOT NULL"
	rows, err := r.pool.Query(ctx, orderSelect+" WHERE "+cond+" ORDER BY o.opened_at DESC", args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []model.ServiceOrder
	for rows.Next() {
		var o model.ServiceOrder
		if err := scanOrder(rows, &o); err != nil {
			return nil, err
		}
		out = append(out, o)
	}
	return out, rows.Err()
}

// NextNumber generates a sequential order number like "OS-2026-000042".
func (r *ServiceOrderRepository) NextNumber(ctx context.Context) (string, error) {
	var seq int64
	if err := r.pool.QueryRow(ctx,
		`SELECT COALESCE(MAX(id), 0) + 1 FROM service_orders`).Scan(&seq); err != nil {
		return "", err
	}
	return fmt.Sprintf("OS-%d-%06d", nowYear(), seq), nil
}
