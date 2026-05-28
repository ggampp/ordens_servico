package handler

import (
	"net/http"

	"github.com/ggampp/ordens_servico/backend/internal/httpx"
	"github.com/ggampp/ordens_servico/backend/internal/service"
)

// DashboardHandler exposes aggregate indicators.
type DashboardHandler struct {
	svc *service.DashboardService
}

// NewDashboardHandler builds a DashboardHandler.
func NewDashboardHandler(svc *service.DashboardService) *DashboardHandler {
	return &DashboardHandler{svc: svc}
}

// Summary handles GET /dashboard.
func (h *DashboardHandler) Summary(w http.ResponseWriter, r *http.Request) {
	s, err := h.svc.Summary(r.Context())
	if err != nil {
		httpx.WriteError(w, err)
		return
	}
	httpx.WriteJSON(w, http.StatusOK, s)
}
