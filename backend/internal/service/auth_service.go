package service

import (
	"context"
	"errors"
	"log/slog"

	"github.com/ggampp/ordens_servico/backend/internal/auth"
	"github.com/ggampp/ordens_servico/backend/internal/httpx"
	"github.com/ggampp/ordens_servico/backend/internal/model"
	"github.com/ggampp/ordens_servico/backend/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

// AuthService handles authentication and user provisioning.
type AuthService struct {
	users *repository.UserRepository
	jwt   *auth.Manager
}

// NewAuthService builds an AuthService.
func NewAuthService(users *repository.UserRepository, jwt *auth.Manager) *AuthService {
	return &AuthService{users: users, jwt: jwt}
}

// Login validates credentials and returns a signed token.
func (s *AuthService) Login(ctx context.Context, in model.LoginInput) (*model.LoginResult, error) {
	u, err := s.users.GetByEmail(ctx, in.Email)
	if errors.Is(err, repository.ErrNotFound) {
		return nil, httpx.NewUnauthorized("invalid credentials")
	}
	if err != nil {
		return nil, err
	}
	if !u.Active {
		return nil, httpx.NewForbidden("user is inactive")
	}
	if bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(in.Password)) != nil {
		return nil, httpx.NewUnauthorized("invalid credentials")
	}
	token, err := s.jwt.Generate(u.ID, u.Role, u.EmployeeID)
	if err != nil {
		return nil, err
	}
	return &model.LoginResult{Token: token, User: *u}, nil
}

// Register creates a new user (admin-only at the handler layer).
func (s *AuthService) Register(ctx context.Context, in model.RegisterInput) (*model.User, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(in.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	u := &model.User{
		Name: in.Name, Email: in.Email, PasswordHash: string(hash),
		Role: in.Role, EmployeeID: in.EmployeeID, Active: true,
	}
	if err := s.users.Create(ctx, u); err != nil {
		if repository.IsUniqueViolation(err) {
			return nil, httpx.NewConflict("email already in use")
		}
		return nil, err
	}
	return u, nil
}

// SeedAdmin creates the initial admin account when no admin exists yet.
func (s *AuthService) SeedAdmin(ctx context.Context, email, password string) error {
	n, err := s.users.CountAdmins(ctx)
	if err != nil {
		return err
	}
	if n > 0 {
		return nil
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	admin := &model.User{
		Name: "Administrator", Email: email, PasswordHash: string(hash),
		Role: model.RoleAdmin, Active: true,
	}
	if err := s.users.Create(ctx, admin); err != nil {
		return err
	}
	slog.Info("seeded initial admin user", "email", email)
	return nil
}
