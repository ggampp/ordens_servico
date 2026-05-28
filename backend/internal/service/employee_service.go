package service

import (
	"context"
	"errors"
	"time"

	"github.com/ggampp/ordens_servico/backend/internal/httpx"
	"github.com/ggampp/ordens_servico/backend/internal/model"
	"github.com/ggampp/ordens_servico/backend/internal/repository"
)

// EmployeeService implements employee business rules.
type EmployeeService struct {
	repo *repository.EmployeeRepository
}

// NewEmployeeService builds an EmployeeService.
func NewEmployeeService(repo *repository.EmployeeRepository) *EmployeeService {
	return &EmployeeService{repo: repo}
}

// Create registers a new employee.
func (s *EmployeeService) Create(ctx context.Context, in model.CreateEmployeeInput) (*model.Employee, error) {
	status := in.Status
	if status == "" {
		status = model.EmployeeActive
	}
	e := &model.Employee{
		Code: in.Code, Name: in.Name, Email: in.Email,
		Phone: in.Phone, Role: in.Role, Status: status,
	}
	if err := s.repo.Create(ctx, e); err != nil {
		if repository.IsUniqueViolation(err) {
			return nil, httpx.NewConflict("employee code already exists")
		}
		return nil, err
	}
	return e, nil
}

// Update modifies an employee.
func (s *EmployeeService) Update(ctx context.Context, id int64, in model.UpdateEmployeeInput) (*model.Employee, error) {
	e, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, mapNotFound(err, "employee not found")
	}
	e.Name = in.Name
	e.Email = in.Email
	e.Phone = in.Phone
	e.Role = in.Role
	if in.Status != "" {
		e.Status = in.Status
	}
	if err := s.repo.Update(ctx, e); err != nil {
		return nil, mapNotFound(err, "employee not found")
	}
	return s.repo.GetByID(ctx, id)
}

// Delete logically removes an employee.
func (s *EmployeeService) Delete(ctx context.Context, id int64) error {
	return mapNotFound(s.repo.SoftDelete(ctx, id), "employee not found")
}

// Get returns a single employee.
func (s *EmployeeService) Get(ctx context.Context, id int64) (*model.Employee, error) {
	e, err := s.repo.GetByID(ctx, id)
	return e, mapNotFound(err, "employee not found")
}

// List returns a paginated set of employees.
func (s *EmployeeService) List(ctx context.Context, f model.EmployeeFilter) (model.PagedResult[model.Employee], error) {
	f.Normalize()
	items, total, err := s.repo.List(ctx, f)
	if err != nil {
		return model.PagedResult[model.Employee]{}, err
	}
	return model.NewPagedResult(items, f.Pagination, total), nil
}

// RecordPosition stores a geolocation point for an employee.
func (s *EmployeeService) RecordPosition(ctx context.Context, id int64, in model.CreatePositionInput) (*model.EmployeePosition, error) {
	ok, err := s.repo.Exists(ctx, id)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, httpx.NewNotFound("employee not found")
	}
	recordedAt := time.Now()
	if in.RecordedAt != nil {
		recordedAt = *in.RecordedAt
	}
	p := &model.EmployeePosition{
		EmployeeID: id, Latitude: in.Latitude,
		Longitude: in.Longitude, RecordedAt: recordedAt,
	}
	if err := s.repo.AddPosition(ctx, p); err != nil {
		return nil, err
	}
	return p, nil
}

// PositionHistory returns paginated location history for an employee.
func (s *EmployeeService) PositionHistory(ctx context.Context, id int64, p model.Pagination) (model.PagedResult[model.EmployeePosition], error) {
	p.Normalize()
	items, total, err := s.repo.PositionHistory(ctx, id, p)
	if err != nil {
		return model.PagedResult[model.EmployeePosition]{}, err
	}
	return model.NewPagedResult(items, p, total), nil
}

// mapNotFound translates repository ErrNotFound into an API 404.
func mapNotFound(err error, msg string) error {
	if err == nil {
		return nil
	}
	if errors.Is(err, repository.ErrNotFound) {
		return httpx.NewNotFound(msg)
	}
	return err
}
