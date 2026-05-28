package handler

import (
	"net/http"
	"strconv"
	"time"

	"github.com/ggampp/ordens_servico/backend/internal/httpx"
	"github.com/ggampp/ordens_servico/backend/internal/middleware"
	"github.com/ggampp/ordens_servico/backend/internal/model"
	"github.com/ggampp/ordens_servico/backend/internal/service"
	"github.com/go-chi/chi/v5"
)

// pathID parses the {id} URL parameter as int64.
func pathID(r *http.Request) (int64, error) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		return 0, httpx.NewBadRequest("invalid id")
	}
	return id, nil
}

// actorFrom derives a service.Actor from the request principal.
func actorFrom(r *http.Request) service.Actor {
	p, _ := middleware.PrincipalFrom(r.Context())
	return service.Actor{UserID: p.UserID, Role: p.Role, EmployeeID: p.EmployeeID}
}

// pagination reads page/page_size query params.
func pagination(r *http.Request) model.Pagination {
	return model.Pagination{
		Page:     httpx.QueryInt(r, "page", 1),
		PageSize: httpx.QueryInt(r, "page_size", 20),
	}
}

// parseDate parses an RFC3339 / date-only query param into a pointer.
func parseDate(r *http.Request, key string) *time.Time {
	v := r.URL.Query().Get(key)
	if v == "" {
		return nil
	}
	for _, layout := range []string{time.RFC3339, "2006-01-02"} {
		if t, err := time.Parse(layout, v); err == nil {
			return &t
		}
	}
	return nil
}

// serviceOrderFilter builds the order filter from query params.
func serviceOrderFilter(r *http.Request) model.ServiceOrderFilter {
	f := model.ServiceOrderFilter{
		Status:     r.URL.Query().Get("status"),
		Priority:   r.URL.Query().Get("priority"),
		DateFrom:   parseDate(r, "date_from"),
		DateTo:     parseDate(r, "date_to"),
		MinLat:     httpx.QueryFloatPtr(r, "min_lat"),
		MinLng:     httpx.QueryFloatPtr(r, "min_lng"),
		MaxLat:     httpx.QueryFloatPtr(r, "max_lat"),
		MaxLng:     httpx.QueryFloatPtr(r, "max_lng"),
		Pagination: pagination(r),
	}
	if v := r.URL.Query().Get("employee_id"); v != "" {
		if id, err := strconv.ParseInt(v, 10, 64); err == nil {
			f.EmployeeID = &id
		}
	}
	return f
}
