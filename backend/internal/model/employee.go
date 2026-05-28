package model

import "time"

// Employee status values.
const (
	EmployeeActive   = "active"
	EmployeeInactive = "inactive"
)

// Employee represents a field worker.
type Employee struct {
	ID        int64     `json:"id"`
	Code      string    `json:"code"`
	Name      string    `json:"name"`
	Email     *string   `json:"email,omitempty"`
	Phone     *string   `json:"phone,omitempty"`
	Role      *string   `json:"role,omitempty"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	// LastPosition is optionally hydrated for map/list views.
	LastPosition *EmployeePosition `json:"last_position,omitempty"`
}

// CreateEmployeeInput is the payload for creating an employee.
type CreateEmployeeInput struct {
	Code   string  `json:"code" validate:"required,max=50"`
	Name   string  `json:"name" validate:"required,max=150"`
	Email  *string `json:"email" validate:"omitempty,email,max=150"`
	Phone  *string `json:"phone" validate:"omitempty,max=30"`
	Role   *string `json:"role" validate:"omitempty,max=80"`
	Status string  `json:"status" validate:"omitempty,oneof=active inactive"`
}

// UpdateEmployeeInput is the payload for updating an employee.
type UpdateEmployeeInput struct {
	Name   string  `json:"name" validate:"required,max=150"`
	Email  *string `json:"email" validate:"omitempty,email,max=150"`
	Phone  *string `json:"phone" validate:"omitempty,max=30"`
	Role   *string `json:"role" validate:"omitempty,max=80"`
	Status string  `json:"status" validate:"omitempty,oneof=active inactive"`
}

// EmployeeFilter holds list query filters.
type EmployeeFilter struct {
	Status string
	Search string
	Pagination
}

// EmployeePosition is a single geolocation record.
type EmployeePosition struct {
	ID         int64     `json:"id"`
	EmployeeID int64     `json:"employee_id"`
	Latitude   float64   `json:"latitude"`
	Longitude  float64   `json:"longitude"`
	RecordedAt time.Time `json:"recorded_at"`
}

// CreatePositionInput is the payload for recording a position.
type CreatePositionInput struct {
	Latitude   float64    `json:"latitude" validate:"required,latitude"`
	Longitude  float64    `json:"longitude" validate:"required,longitude"`
	RecordedAt *time.Time `json:"recorded_at"`
}
