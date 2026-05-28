package model

import "time"

// Service order status values.
const (
	StatusOpen       = "open"
	StatusAssigned   = "assigned"
	StatusInProgress = "in_progress"
	StatusCompleted  = "completed"
	StatusCancelled  = "cancelled"
)

// Priority values.
const (
	PriorityLow    = "low"
	PriorityMedium = "medium"
	PriorityHigh   = "high"
	PriorityUrgent = "urgent"
)

// ValidStatusTransitions defines the allowed status state machine.
var ValidStatusTransitions = map[string][]string{
	StatusOpen:       {StatusAssigned, StatusInProgress, StatusCancelled},
	StatusAssigned:   {StatusInProgress, StatusCancelled, StatusOpen},
	StatusInProgress: {StatusCompleted, StatusCancelled},
	StatusCompleted:  {},
	StatusCancelled:  {},
}

// CanTransition reports whether moving from->to is allowed.
func CanTransition(from, to string) bool {
	if from == to {
		return true
	}
	for _, s := range ValidStatusTransitions[from] {
		if s == to {
			return true
		}
	}
	return false
}

// ServiceOrder represents a work order.
type ServiceOrder struct {
	ID           int64      `json:"id"`
	Number       string     `json:"number"`
	Title        string     `json:"title"`
	Description  *string    `json:"description,omitempty"`
	Priority     string     `json:"priority"`
	Status       string     `json:"status"`
	EmployeeID   *int64     `json:"employee_id,omitempty"`
	EmployeeName *string    `json:"employee_name,omitempty"`
	Address      *string    `json:"address,omitempty"`
	Latitude     *float64   `json:"latitude,omitempty"`
	Longitude    *float64   `json:"longitude,omitempty"`
	OpenedAt     time.Time  `json:"opened_at"`
	DueAt        *time.Time `json:"due_at,omitempty"`
	CompletedAt  *time.Time `json:"completed_at,omitempty"`
	Notes        *string    `json:"notes,omitempty"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
}

// CreateServiceOrderInput is the payload for creating an order.
type CreateServiceOrderInput struct {
	Number      string     `json:"number" validate:"omitempty,max=30"`
	Title       string     `json:"title" validate:"required,max=200"`
	Description *string    `json:"description"`
	Priority    string     `json:"priority" validate:"omitempty,oneof=low medium high urgent"`
	EmployeeID  *int64     `json:"employee_id"`
	Address     *string    `json:"address" validate:"omitempty,max=300"`
	Latitude    *float64   `json:"latitude" validate:"omitempty,latitude"`
	Longitude   *float64   `json:"longitude" validate:"omitempty,longitude"`
	DueAt       *time.Time `json:"due_at"`
	Notes       *string    `json:"notes"`
}

// UpdateServiceOrderInput is the payload for editing an order.
type UpdateServiceOrderInput struct {
	Title       string     `json:"title" validate:"required,max=200"`
	Description *string    `json:"description"`
	Priority    string     `json:"priority" validate:"omitempty,oneof=low medium high urgent"`
	Address     *string    `json:"address" validate:"omitempty,max=300"`
	Latitude    *float64   `json:"latitude" validate:"omitempty,latitude"`
	Longitude   *float64   `json:"longitude" validate:"omitempty,longitude"`
	DueAt       *time.Time `json:"due_at"`
	Notes       *string    `json:"notes"`
}

// StatusChangeInput changes an order's status.
type StatusChangeInput struct {
	Status string `json:"status" validate:"required,oneof=open assigned in_progress completed cancelled"`
	Note   string `json:"note" validate:"omitempty,max=500"`
}

// AssignInput assigns an order to an employee.
type AssignInput struct {
	EmployeeID int64 `json:"employee_id" validate:"required"`
}

// ServiceOrderFilter holds list query filters.
type ServiceOrderFilter struct {
	Status     string
	Priority   string
	EmployeeID *int64
	DateFrom   *time.Time
	DateTo     *time.Time
	// Bounding box for geographic region filtering.
	MinLat, MinLng, MaxLat, MaxLng *float64
	Pagination
}

// ServiceOrderHistory is a status change record.
type ServiceOrderHistory struct {
	ID             int64     `json:"id"`
	ServiceOrderID int64     `json:"service_order_id"`
	OldStatus      *string   `json:"old_status,omitempty"`
	NewStatus      string    `json:"new_status"`
	ChangedBy      *int64    `json:"changed_by,omitempty"`
	ChangedByName  *string   `json:"changed_by_name,omitempty"`
	Note           *string   `json:"note,omitempty"`
	ChangedAt      time.Time `json:"changed_at"`
}
