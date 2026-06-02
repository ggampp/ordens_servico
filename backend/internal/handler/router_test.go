package handler

import (
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/ggampp/ordens_servico/backend/internal/auth"
)

// newTestRouter builds a router with empty handlers; the SPA handler does not
// touch the database, so this is enough to exercise static serving and routing.
func newTestRouter(t *testing.T, staticDir string) http.Handler {
	t.Helper()
	jwt := auth.NewManager("test-secret", time.Hour)
	return NewRouter(Handlers{}, jwt, staticDir)
}

func TestSPAServesIndexFallback(t *testing.T) {
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "index.html"), []byte("<!doctype html><title>spa</title>"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "app.js"), []byte("console.log('hi')"), 0o644); err != nil {
		t.Fatal(err)
	}

	r := newTestRouter(t, dir)

	cases := []struct {
		path     string
		wantBody string
	}{
		{"/", "<!doctype html>"},        // root -> index
		{"/login", "<!doctype html>"},   // client route -> index fallback
		{"/app.js", "console.log"},      // existing asset -> served as-is
	}
	for _, c := range cases {
		req := httptest.NewRequest(http.MethodGet, c.path, nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		if w.Code != http.StatusOK {
			t.Fatalf("%s: status = %d, want 200", c.path, w.Code)
		}
		if body := w.Body.String(); !contains(body, c.wantBody) {
			t.Fatalf("%s: body %q does not contain %q", c.path, body, c.wantBody)
		}
	}
}

func TestSPADoesNotSwallowUnknownAPIPaths(t *testing.T) {
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "index.html"), []byte("<!doctype html>"), 0o644); err != nil {
		t.Fatal(err)
	}

	r := newTestRouter(t, dir)
	req := httptest.NewRequest(http.MethodGet, "/api/v1/does-not-exist", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Fatalf("status = %d, want 404", w.Code)
	}
	if ct := w.Header().Get("Content-Type"); !contains(ct, "application/json") {
		t.Fatalf("content-type = %q, want application/json", ct)
	}
}

func TestNoStaticDirLeavesDefault404(t *testing.T) {
	r := newTestRouter(t, "")
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusNotFound {
		t.Fatalf("status = %d, want 404", w.Code)
	}
}

func TestHealthEndpoint(t *testing.T) {
	r := newTestRouter(t, "")
	for _, path := range []string{"/health", "/healthz"} {
		req := httptest.NewRequest(http.MethodGet, path, nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		if w.Code != http.StatusOK {
			t.Fatalf("%s: status = %d, want 200", path, w.Code)
		}
	}
}

func contains(s, sub string) bool {
	return len(s) >= len(sub) && (s == sub || indexOf(s, sub) >= 0)
}

func indexOf(s, sub string) int {
	for i := 0; i+len(sub) <= len(s); i++ {
		if s[i:i+len(sub)] == sub {
			return i
		}
	}
	return -1
}
