package handler

import (
	"net/http"

	"github.com/ggampp/ordens_servico/backend/internal/httpx"
	"github.com/ggampp/ordens_servico/backend/internal/model"
	"github.com/ggampp/ordens_servico/backend/internal/service"
)

// EmployeeHandler exposes employee CRUD and position endpoints.
type EmployeeHandler struct {
	svc *service.EmployeeService
}

// NewEmployeeHandler builds an EmployeeHandler.
func NewEmployeeHandler(svc *service.EmployeeService) *EmployeeHandler {
	return &EmployeeHandler{svc: svc}
}

// Create handles POST /employees.
func (h *EmployeeHandler) Create(w http.ResponseWriter, r *http.Request) {
	var in model.CreateEmployeeInput
	if err := httpx.DecodeAndValidate(r, &in); err != nil {
		httpx.WriteError(w, err)
		return
	}
	e, err := h.svc.Create(r.Context(), in)
	if err != nil {
		httpx.WriteError(w, err)
		return
	}
	httpx.WriteJSON(w, http.StatusCreated, e)
}

// List handles GET /employees.
func (h *EmployeeHandler) List(w http.ResponseWriter, r *http.Request) {
	f := model.EmployeeFilter{
		Status:     r.URL.Query().Get("status"),
		Search:     r.URL.Query().Get("search"),
		Pagination: pagination(r),
	}
	res, err := h.svc.List(r.Context(), f)
	if err != nil {
		httpx.WriteError(w, err)
		return
	}
	httpx.WriteJSON(w, http.StatusOK, res)
}

// Get handles GET /employees/{id}.
func (h *EmployeeHandler) Get(w http.ResponseWriter, r *http.Request) {
	id, err := pathID(r)
	if err != nil {
		httpx.WriteError(w, err)
		return
	}
	e, err := h.svc.Get(r.Context(), id)
	if err != nil {
		httpx.WriteError(w, err)
		return
	}
	httpx.WriteJSON(w, http.StatusOK, e)
}

// Update handles PUT /employees/{id}.
func (h *EmployeeHandler) Update(w http.ResponseWriter, r *http.Request) {
	id, err := pathID(r)
	if err != nil {
		httpx.WriteError(w, err)
		return
	}
	var in model.UpdateEmployeeInput
	if err := httpx.DecodeAndValidate(r, &in); err != nil {
		httpx.WriteError(w, err)
		return
	}
	e, err := h.svc.Update(r.Context(), id, in)
	if err != nil {
		httpx.WriteError(w, err)
		return
	}
	httpx.WriteJSON(w, http.StatusOK, e)
}

// Delete handles DELETE /employees/{id}.
func (h *EmployeeHandler) Delete(w http.ResponseWriter, r *http.Request) {
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

// RecordPosition handles POST /employees/{id}/position.
func (h *EmployeeHandler) RecordPosition(w http.ResponseWriter, r *http.Request) {
	id, err := pathID(r)
	if err != nil {
		httpx.WriteError(w, err)
		return
	}
	var in model.CreatePositionInput
	if err := httpx.DecodeAndValidate(r, &in); err != nil {
		httpx.WriteError(w, err)
		return
	}
	p, err := h.svc.RecordPosition(r.Context(), id, in)
	if err != nil {
		httpx.WriteError(w, err)
		return
	}
	httpx.WriteJSON(w, http.StatusCreated, p)
}

// PositionHistory handles GET /employees/{id}/positions.
func (h *EmployeeHandler) PositionHistory(w http.ResponseWriter, r *http.Request) {
	id, err := pathID(r)
	if err != nil {
		httpx.WriteError(w, err)
		return
	}
	res, err := h.svc.PositionHistory(r.Context(), id, pagination(r))
	if err != nil {
		httpx.WriteError(w, err)
		return
	}
	httpx.WriteJSON(w, http.StatusOK, res)
}
