package handler

import (
	"net/http"

	"github.com/ggampp/ordens_servico/backend/internal/httpx"
	"github.com/ggampp/ordens_servico/backend/internal/service"
)

// MapHandler exposes geolocation views.
type MapHandler struct {
	svc *service.MapService
}

// NewMapHandler builds a MapHandler.
func NewMapHandler(svc *service.MapService) *MapHandler {
	return &MapHandler{svc: svc}
}

// Employees handles GET /map/employees.
func (h *MapHandler) Employees(w http.ResponseWriter, r *http.Request) {
	emps, err := h.svc.Employees(r.Context())
	if err != nil {
		httpx.WriteError(w, err)
		return
	}
	httpx.WriteJSON(w, http.StatusOK, emps)
}

// ServiceOrders handles GET /map/service-orders.
func (h *MapHandler) ServiceOrders(w http.ResponseWriter, r *http.Request) {
	orders, err := h.svc.ServiceOrders(r.Context(), serviceOrderFilter(r), actorFrom(r))
	if err != nil {
		httpx.WriteError(w, err)
		return
	}
	httpx.WriteJSON(w, http.StatusOK, orders)
}

// Overview handles GET /map/overview.
func (h *MapHandler) Overview(w http.ResponseWriter, r *http.Request) {
	res, err := h.svc.Overview(r.Context(), serviceOrderFilter(r), actorFrom(r))
	if err != nil {
		httpx.WriteError(w, err)
		return
	}
	httpx.WriteJSON(w, http.StatusOK, res)
}
