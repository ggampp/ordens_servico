package model

import "time"

// User roles.
const (
	RoleAdmin      = "admin"
	RoleSupervisor = "supervisor"
	RoleOperator   = "operator"
)

// User is an authentication principal.
type User struct {
	ID           int64     `json:"id"`
	Name         string    `json:"name"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"-"`
	Role         string    `json:"role"`
	EmployeeID   *int64    `json:"employee_id,omitempty"`
	Active       bool      `json:"active"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// LoginInput is the login payload.
type LoginInput struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

// RegisterInput creates a user (admin only).
type RegisterInput struct {
	Name       string `json:"name" validate:"required,max=150"`
	Email      string `json:"email" validate:"required,email,max=150"`
	Password   string `json:"password" validate:"required,min=6,max=72"`
	Role       string `json:"role" validate:"required,oneof=admin supervisor operator"`
	EmployeeID *int64 `json:"employee_id"`
}

// LoginResult is returned on successful authentication.
type LoginResult struct {
	Token string `json:"token"`
	User  User   `json:"user"`
}
