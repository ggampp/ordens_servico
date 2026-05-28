package handler

import (
	"net/http"
	"time"

	"github.com/ggampp/ordens_servico/backend/internal/auth"
	"github.com/ggampp/ordens_servico/backend/internal/httpx"
	mw "github.com/ggampp/ordens_servico/backend/internal/middleware"
	"github.com/ggampp/ordens_servico/backend/internal/model"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
)

// Handlers bundles all HTTP handlers for route registration.
type Handlers struct {
	Auth      *AuthHandler
	Employee  *EmployeeHandler
	Order     *ServiceOrderHandler
	Map       *MapHandler
	Dashboard *DashboardHandler
}

// NewRouter wires middleware, routes and role-based access control.
func NewRouter(h Handlers, jwt *auth.Manager) http.Handler {
	r := chi.NewRouter()

	r.Use(mw.Recoverer)
	r.Use(mw.RequestLogger)
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	// Health check.
	r.Get("/health", func(w http.ResponseWriter, _ *http.Request) {
		httpx.WriteJSON(w, http.StatusOK, map[string]any{"status": "ok", "time": time.Now()})
	})

	// API documentation (Swagger UI + OpenAPI spec).
	r.Get("/swagger", swaggerUI)
	r.Get("/swagger/", swaggerUI)
	r.Get("/openapi.yaml", openAPISpec)

	r.Route("/api/v1", func(r chi.Router) {
		// Public auth endpoints.
		r.Post("/auth/login", h.Auth.Login)

		// Authenticated endpoints.
		r.Group(func(r chi.Router) {
			r.Use(mw.Authenticator(jwt))

			// User provisioning (admin only).
			r.With(mw.RequireRole(model.RoleAdmin)).
				Post("/auth/register", h.Auth.Register)

			// ----- Employees -----
			r.Route("/employees", func(r chi.Router) {
				r.Get("/", h.Employee.List)
				r.Get("/{id}", h.Employee.Get)
				r.Get("/{id}/positions", h.Employee.PositionHistory)
				// Any authenticated user may report a position.
				r.Post("/{id}/position", h.Employee.RecordPosition)

				r.Group(func(r chi.Router) {
					r.Use(mw.RequireRole(model.RoleAdmin, model.RoleSupervisor))
					r.Post("/", h.Employee.Create)
					r.Put("/{id}", h.Employee.Update)
					r.Delete("/{id}", h.Employee.Delete)
				})
			})

			// ----- Service Orders -----
			r.Route("/service-orders", func(r chi.Router) {
				r.Get("/", h.Order.List)
				r.Get("/{id}", h.Order.Get)
				r.Get("/{id}/history", h.Order.History)
				// Operators may update and change status of their own orders.
				r.Put("/{id}", h.Order.Update)
				r.Patch("/{id}/status", h.Order.ChangeStatus)

				r.Group(func(r chi.Router) {
					r.Use(mw.RequireRole(model.RoleAdmin, model.RoleSupervisor))
					r.Post("/", h.Order.Create)
					r.Delete("/{id}", h.Order.Delete)
					r.Patch("/{id}/assign", h.Order.Assign)
				})
			})

			// ----- Map -----
			r.Route("/map", func(r chi.Router) {
				r.Get("/employees", h.Map.Employees)
				r.Get("/service-orders", h.Map.ServiceOrders)
				r.Get("/overview", h.Map.Overview)
			})

			// ----- Dashboard (admin/supervisor) -----
			r.With(mw.RequireRole(model.RoleAdmin, model.RoleSupervisor)).
				Get("/dashboard", h.Dashboard.Summary)
		})
	})

	return r
}
