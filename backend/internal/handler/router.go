package handler

import (
	"net/http"
	"os"
	"path/filepath"
	"strings"
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
//
// When staticDir points to a directory containing the built SPA, the backend
// also serves the frontend, turning the service into a single-port monolith
// (API under /api/v1, SPA everywhere else). Pass "" to disable.
func NewRouter(h Handlers, jwt *auth.Manager, staticDir string) http.Handler {
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

	// Serve the built SPA (single-port monolith) when a static dir is provided.
	if spa := newSPAHandler(staticDir); spa != nil {
		r.NotFound(spa)
	}

	return r
}

// newSPAHandler returns a handler that serves static assets from dir and falls
// back to index.html for client-side routes. It returns nil when dir is empty
// or has no index.html, leaving chi's default 404 in place (e.g. in tests).
func newSPAHandler(dir string) http.HandlerFunc {
	if dir == "" {
		return nil
	}
	index := filepath.Join(dir, "index.html")
	if _, err := os.Stat(index); err != nil {
		return nil
	}
	fs := http.Dir(dir)

	return func(w http.ResponseWriter, req *http.Request) {
		// Never let unknown API paths fall through to the SPA.
		if strings.HasPrefix(req.URL.Path, "/api/") {
			httpx.WriteError(w, httpx.NewNotFound("recurso não encontrado"))
			return
		}

		// Serve the requested asset when it exists; otherwise the SPA shell.
		clean := filepath.Clean(req.URL.Path)
		if clean != "/" && clean != "." {
			if f, err := fs.Open(clean); err == nil {
				if info, statErr := f.Stat(); statErr == nil && !info.IsDir() {
					http.ServeContent(w, req, info.Name(), info.ModTime(), f)
					_ = f.Close()
					return
				}
				_ = f.Close()
			}
		}
		http.ServeFile(w, req, index)
	}
}
