package handler

import (
	"net/http"

	"github.com/ggampp/ordens_servico/backend/internal/httpx"
	"github.com/ggampp/ordens_servico/backend/internal/model"
	"github.com/ggampp/ordens_servico/backend/internal/service"
)

// AuthHandler exposes authentication endpoints.
type AuthHandler struct {
	svc *service.AuthService
}

// NewAuthHandler builds an AuthHandler.
func NewAuthHandler(svc *service.AuthService) *AuthHandler {
	return &AuthHandler{svc: svc}
}

// Login godoc
// @Summary Authenticate and obtain a JWT
// @Tags auth
// @Accept json
// @Produce json
// @Param body body model.LoginInput true "Credentials"
// @Success 200 {object} model.LoginResult
// @Router /auth/login [post]
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var in model.LoginInput
	if err := httpx.DecodeAndValidate(r, &in); err != nil {
		httpx.WriteError(w, err)
		return
	}
	res, err := h.svc.Login(r.Context(), in)
	if err != nil {
		httpx.WriteError(w, err)
		return
	}
	httpx.WriteJSON(w, http.StatusOK, res)
}

// Register godoc
// @Summary Create a new user (admin only)
// @Tags auth
// @Security BearerAuth
// @Param body body model.RegisterInput true "User"
// @Success 201 {object} model.User
// @Router /auth/register [post]
func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var in model.RegisterInput
	if err := httpx.DecodeAndValidate(r, &in); err != nil {
		httpx.WriteError(w, err)
		return
	}
	u, err := h.svc.Register(r.Context(), in)
	if err != nil {
		httpx.WriteError(w, err)
		return
	}
	httpx.WriteJSON(w, http.StatusCreated, u)
}
