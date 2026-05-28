package service

import (
	"context"

	"github.com/ggampp/ordens_servico/backend/internal/httpx"
	"github.com/ggampp/ordens_servico/backend/internal/model"
	"github.com/ggampp/ordens_servico/backend/internal/repository"
)

// Actor identifies the authenticated principal acting on a resource.
type Actor struct {
	UserID     int64
	Role       string
	EmployeeID *int64
}

func (a Actor) isOperator() bool {
	return a.Role == model.RoleOperator
}

// ServiceOrderService implements service order business rules.
type ServiceOrderService struct {
	repo      *repository.ServiceOrderRepository
	employees *repository.EmployeeRepository
}

// NewServiceOrderService builds a ServiceOrderService.
func NewServiceOrderService(repo *repository.ServiceOrderRepository, employees *repository.EmployeeRepository) *ServiceOrderService {
	return &ServiceOrderService{repo: repo, employees: employees}
}

// Create registers a new order, auto-generating the number when omitted.
func (s *ServiceOrderService) Create(ctx context.Context, in model.CreateServiceOrderInput, actor Actor) (*model.ServiceOrder, error) {
	number := in.Number
	if number == "" {
		n, err := s.repo.NextNumber(ctx)
		if err != nil {
			return nil, err
		}
		number = n
	}
	priority := in.Priority
	if priority == "" {
		priority = model.PriorityMedium
	}
	status := model.StatusOpen
	if in.EmployeeID != nil {
		if err := s.assertEmployee(ctx, *in.EmployeeID); err != nil {
			return nil, err
		}
		status = model.StatusAssigned
	}
	o := &model.ServiceOrder{
		Number: number, Title: in.Title, Description: in.Description,
		Priority: priority, Status: status, EmployeeID: in.EmployeeID,
		Address: in.Address, Latitude: in.Latitude, Longitude: in.Longitude,
		DueAt: in.DueAt, Notes: in.Notes,
	}
	if err := s.repo.Create(ctx, o, &actor.UserID); err != nil {
		if repository.IsUniqueViolation(err) {
			return nil, httpx.NewConflict("order number already exists")
		}
		return nil, err
	}
	return s.repo.GetByID(ctx, o.ID)
}

// Update edits an order, enforcing operator scope.
func (s *ServiceOrderService) Update(ctx context.Context, id int64, in model.UpdateServiceOrderInput, actor Actor) (*model.ServiceOrder, error) {
	existing, err := s.getScoped(ctx, id, actor)
	if err != nil {
		return nil, err
	}
	existing.Title = in.Title
	existing.Description = in.Description
	if in.Priority != "" {
		existing.Priority = in.Priority
	}
	existing.Address = in.Address
	existing.Latitude = in.Latitude
	existing.Longitude = in.Longitude
	existing.DueAt = in.DueAt
	existing.Notes = in.Notes
	if err := s.repo.Update(ctx, existing); err != nil {
		return nil, mapNotFound(err, "service order not found")
	}
	return s.repo.GetByID(ctx, id)
}

// Get returns a single order, enforcing operator scope.
func (s *ServiceOrderService) Get(ctx context.Context, id int64, actor Actor) (*model.ServiceOrder, error) {
	return s.getScoped(ctx, id, actor)
}

// List returns filtered, paginated orders. Operators are restricted to their
// own assignments regardless of the requested employee filter.
func (s *ServiceOrderService) List(ctx context.Context, f model.ServiceOrderFilter, actor Actor) (model.PagedResult[model.ServiceOrder], error) {
	f.Normalize()
	if actor.isOperator() {
		f.EmployeeID = actor.EmployeeID
	}
	items, total, err := s.repo.List(ctx, f)
	if err != nil {
		return model.PagedResult[model.ServiceOrder]{}, err
	}
	return model.NewPagedResult(items, f.Pagination, total), nil
}

// Delete logically removes an order (supervisors/admins only at handler layer).
func (s *ServiceOrderService) Delete(ctx context.Context, id int64) error {
	return mapNotFound(s.repo.SoftDelete(ctx, id), "service order not found")
}

// ChangeStatus transitions an order's status, validating the state machine.
func (s *ServiceOrderService) ChangeStatus(ctx context.Context, id int64, in model.StatusChangeInput, actor Actor) (*model.ServiceOrder, error) {
	existing, err := s.getScoped(ctx, id, actor)
	if err != nil {
		return nil, err
	}
	if !model.CanTransition(existing.Status, in.Status) {
		return nil, httpx.NewBadRequest("invalid status transition from " + existing.Status + " to " + in.Status)
	}
	if err := s.repo.ChangeStatus(ctx, id, existing.Status, in.Status, &actor.UserID, in.Note); err != nil {
		return nil, mapNotFound(err, "service order not found")
	}
	return s.repo.GetByID(ctx, id)
}

// Assign sets the responsible employee (supervisors/admins only).
func (s *ServiceOrderService) Assign(ctx context.Context, id int64, in model.AssignInput, actor Actor) (*model.ServiceOrder, error) {
	existing, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, mapNotFound(err, "service order not found")
	}
	if err := s.assertEmployee(ctx, in.EmployeeID); err != nil {
		return nil, err
	}
	if err := s.repo.Assign(ctx, id, in.EmployeeID, existing.Status, &actor.UserID); err != nil {
		return nil, mapNotFound(err, "service order not found")
	}
	return s.repo.GetByID(ctx, id)
}

// History returns the status change history, enforcing operator scope.
func (s *ServiceOrderService) History(ctx context.Context, id int64, actor Actor) ([]model.ServiceOrderHistory, error) {
	if _, err := s.getScoped(ctx, id, actor); err != nil {
		return nil, err
	}
	return s.repo.History(ctx, id)
}

// getScoped fetches an order and enforces that operators only access their own.
func (s *ServiceOrderService) getScoped(ctx context.Context, id int64, actor Actor) (*model.ServiceOrder, error) {
	o, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, mapNotFound(err, "service order not found")
	}
	if actor.isOperator() {
		if actor.EmployeeID == nil || o.EmployeeID == nil || *o.EmployeeID != *actor.EmployeeID {
			return nil, httpx.NewForbidden("you can only access orders assigned to you")
		}
	}
	return o, nil
}

func (s *ServiceOrderService) assertEmployee(ctx context.Context, employeeID int64) error {
	ok, err := s.employees.Exists(ctx, employeeID)
	if err != nil {
		return err
	}
	if !ok {
		return httpx.NewBadRequest("assigned employee does not exist")
	}
	return nil
}
