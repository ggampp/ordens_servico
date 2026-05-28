package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/ggampp/ordens_servico/backend/internal/auth"
	"github.com/ggampp/ordens_servico/backend/internal/httpx"
)

type ctxKey string

const principalKey ctxKey = "principal"

// Principal is the authenticated user attached to the request context.
type Principal struct {
	UserID     int64
	Role       string
	EmployeeID *int64
}

// Authenticator validates JWTs and injects the principal.
func Authenticator(jwt *auth.Manager) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			header := r.Header.Get("Authorization")
			if !strings.HasPrefix(header, "Bearer ") {
				httpx.WriteError(w, httpx.NewUnauthorized("missing bearer token"))
				return
			}
			claims, err := jwt.Parse(strings.TrimPrefix(header, "Bearer "))
			if err != nil {
				httpx.WriteError(w, httpx.NewUnauthorized("invalid or expired token"))
				return
			}
			p := Principal{UserID: claims.UserID, Role: claims.Role, EmployeeID: claims.EmployeeID}
			ctx := context.WithValue(r.Context(), principalKey, p)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// RequireRole restricts a route to one of the allowed roles.
func RequireRole(roles ...string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			p, ok := PrincipalFrom(r.Context())
			if !ok {
				httpx.WriteError(w, httpx.NewUnauthorized("authentication required"))
				return
			}
			for _, role := range roles {
				if p.Role == role {
					next.ServeHTTP(w, r)
					return
				}
			}
			httpx.WriteError(w, httpx.NewForbidden("insufficient permissions"))
		})
	}
}

// PrincipalFrom extracts the authenticated principal from a context.
func PrincipalFrom(ctx context.Context) (Principal, bool) {
	p, ok := ctx.Value(principalKey).(Principal)
	return p, ok
}
