package handler

import (
	"net/http"

	"github.com/ggampp/ordens_servico/backend/internal/httpx"
	"github.com/ggampp/ordens_servico/backend/internal/model"
	"github.com/ggampp/ordens_servico/backend/internal/service"
)

// ServiceOrderHandler exposes service order endpoints.
type ServiceOrderHandler struct {
	svc *service.ServiceOrderService
}

// NewServiceOrderHandler builds a ServiceOrderHandler.
func NewServiceOrderHandler(svc *service.ServiceOrderService) *ServiceOrderHandler {
	return &ServiceOrderHandler{svc: svc}
}

// Create handles POST /service-orders.
func (h *ServiceOrderHandler) Create(w http.ResponseWriter, r *http.Request) {
	var in model.CreateServiceOrderInput
	if err := httpx.DecodeAndValidate(r, &in); err != nil {
		httpx.WriteError(w, err)
		return
	}
	o, err := h.svc.Create(r.Context(), in, actorFrom(r))
	if err != nil {
		httpx.WriteError(w, err)
		return
	}
	httpx.WriteJSON(w, http.StatusCreated, o)
}

// List handles GET /service-orders.
func (h *ServiceOrderHandler) List(w http.ResponseWriter, r *http.Request) {
	res, err := h.svc.List(r.Context(), serviceOrderFilter(r), actorFrom(r))
	if err != nil {
		httpx.WriteError(w, err)
		return
	}
	httpx.WriteJSON(w, http.StatusOK, res)
}

// Get handles GET /service-orders/{id}.
func (h *ServiceOrderHandler) Get(w http.ResponseWriter, r *http.Request) {
	id, err := pathID(r)
	if err != nil {
		httpx.WriteError(w, err)
		return
	}
	o, err := h.svc.Get(r.Context(), id, actorFrom(r))
	if err != nil {
		httpx.WriteError(w, err)
		return
	}
	httpx.WriteJSON(w, http.StatusOK, o)
}

// Update handles PUT /service-orders/{id}.
func (h *ServiceOrderHandler) Update(w http.ResponseWriter, r *http.Request) {
	id, err := pathID(r)
	if err != nil {
		httpx.WriteError(w, err)
		return
	}
	var in model.UpdateServiceOrderInput
	if err := httpx.DecodeAndValidate(r, &in); err != nil {
		httpx.WriteError(w, err)
		return
	}
	o, err := h.svc.Update(r.Context(), id, in, actorFrom(r))
	if err != nil {
		httpx.WriteError(w, err)
		return
	}
	httpx.WriteJSON(w, http.StatusOK, o)
}

// Delete handles DELETE /service-orders/{id}.
func (h *ServiceOrderHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := pathID(r)
	if err != nil {
		httpx.WriteError(w, err)
		return
	}
	if err := h.svc.Delete(r.Context(), id); err != nil {
		httpx.WriteError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// ChangeStatus handles PATCH /service-orders/{id}/status.
func (h *ServiceOrderHandler) ChangeStatus(w http.ResponseWriter, r *http.Request) {
	id, err := pathID(r)
	if err != nil {
		httpx.WriteError(w, err)
		return
	}
	var in model.StatusChangeInput
	if err := httpx.DecodeAndValidate(r, &in); err != nil {
		httpx.WriteError(w, err)
		return
	}
	o, err := h.svc.ChangeStatus(r.Context(), id, in, actorFrom(r))
	if err != nil {
		httpx.WriteError(w, err)
		return
	}
	httpx.WriteJSON(w, http.StatusOK, o)
}

// Assign handles PATCH /service-orders/{id}/assign.
func (h *ServiceOrderHandler) Assign(w http.ResponseWriter, r *http.Request) {
	id, err := pathID(r)
	if err != nil {
		httpx.WriteError(w, err)
		return
	}
	var in model.AssignInput
	if err := httpx.DecodeAndValidate(r, &in); err != nil {
		httpx.WriteError(w, err)
		return
	}
	o, err := h.svc.Assign(r.Context(), id, in, actorFrom(r))
	if err != nil {
		httpx.WriteError(w, err)
		return
	}
	httpx.WriteJSON(w, http.StatusOK, o)
}

// History handles GET /service-orders/{id}/history.
func (h *ServiceOrderHandler) History(w http.ResponseWriter, r *http.Request) {
	id, err := pathID(r)
	if err != nil {
		httpx.WriteError(w, err)
		return
	}
	hist, err := h.svc.History(r.Context(), id, actorFrom(r))
	if err != nil {
		httpx.WriteError(w, err)
		return
	}
	if hist == nil {
		hist = []model.ServiceOrderHistory{}
	}
	httpx.WriteJSON(w, http.StatusOK, hist)
}
