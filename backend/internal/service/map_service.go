package service

import (
	"context"

	"github.com/ggampp/ordens_servico/backend/internal/model"
	"github.com/ggampp/ordens_servico/backend/internal/repository"
)

// MapService assembles geolocation views.
type MapService struct {
	orders  *repository.ServiceOrderRepository
	mapRepo *repository.MapRepository
}

// NewMapService builds a MapService.
func NewMapService(orders *repository.ServiceOrderRepository, mapRepo *repository.MapRepository) *MapService {
	return &MapService{orders: orders, mapRepo: mapRepo}
}

// MapOverview bundles employees and orders for a single map load.
type MapOverview struct {
	Employees     []model.Employee     `json:"employees"`
	ServiceOrders []model.ServiceOrder `json:"service_orders"`
}

// Employees returns active employees with their latest position.
func (s *MapService) Employees(ctx context.Context) ([]model.Employee, error) {
	emps, err := s.mapRepo.EmployeesWithPosition(ctx)
	if err != nil {
		return nil, err
	}
	if emps == nil {
		emps = []model.Employee{}
	}
	return emps, nil
}

// ServiceOrders returns geolocated orders matching the filter.
func (s *MapService) ServiceOrders(ctx context.Context, f model.ServiceOrderFilter, actor Actor) ([]model.ServiceOrder, error) {
	if actor.isOperator() {
		f.EmployeeID = actor.EmployeeID
	}
	orders, err := s.orders.ForMap(ctx, f)
	if err != nil {
		return nil, err
	}
	if orders == nil {
		orders = []model.ServiceOrder{}
	}
	return orders, nil
}

// Overview returns employees and orders together.
func (s *MapService) Overview(ctx context.Context, f model.ServiceOrderFilter, actor Actor) (*MapOverview, error) {
	emps, err := s.Employees(ctx)
	if err != nil {
		return nil, err
	}
	orders, err := s.ServiceOrders(ctx, f, actor)
	if err != nil {
		return nil, err
	}
	return &MapOverview{Employees: emps, ServiceOrders: orders}, nil
}
