package httpx

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
)

// APIError is a domain error carrying an HTTP status code.
type APIError struct {
	Status  int    `json:"-"`
	Code    string `json:"code"`
	Message string `json:"message"`
}

func (e *APIError) Error() string { return e.Message }

// Constructors for common error categories.
func NewBadRequest(msg string) *APIError { return &APIError{http.StatusBadRequest, "bad_request", msg} }
func NewUnauthorized(msg string) *APIError {
	return &APIError{http.StatusUnauthorized, "unauthorized", msg}
}
func NewForbidden(msg string) *APIError { return &APIError{http.StatusForbidden, "forbidden", msg} }
func NewNotFound(msg string) *APIError  { return &APIError{http.StatusNotFound, "not_found", msg} }
func NewConflict(msg string) *APIError  { return &APIError{http.StatusConflict, "conflict", msg} }
func NewInternal(msg string) *APIError {
	return &APIError{http.StatusInternalServerError, "internal", msg}
}

// errorResponse is the JSON body returned for failures.
type errorResponse struct {
	Error APIError `json:"error"`
}

// WriteJSON serializes v with the given status code.
func WriteJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if v != nil {
		_ = json.NewEncoder(w).Encode(v)
	}
}

// WriteError centralizes error rendering. Unknown errors are logged and
// surfaced as a generic 500 so internals never leak to clients.
func WriteError(w http.ResponseWriter, err error) {
	var apiErr *APIError
	if errors.As(err, &apiErr) {
		WriteJSON(w, apiErr.Status, errorResponse{Error: *apiErr})
		return
	}
	slog.Error("unhandled error", "error", err)
	generic := NewInternal("internal server error")
	WriteJSON(w, generic.Status, errorResponse{Error: *generic})
}
