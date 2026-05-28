package auth

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// Claims is the JWT payload.
type Claims struct {
	UserID     int64  `json:"uid"`
	Role       string `json:"role"`
	EmployeeID *int64 `json:"eid,omitempty"`
	jwt.RegisteredClaims
}

// Manager issues and validates JWT tokens.
type Manager struct {
	secret []byte
	expiry time.Duration
}

// NewManager builds a JWT manager.
func NewManager(secret string, expiry time.Duration) *Manager {
	return &Manager{secret: []byte(secret), expiry: expiry}
}

// Generate creates a signed token for the given principal.
func (m *Manager) Generate(userID int64, role string, employeeID *int64) (string, error) {
	now := time.Now()
	claims := Claims{
		UserID:     userID,
		Role:       role,
		EmployeeID: employeeID,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(m.expiry)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(m.secret)
}

// Parse validates a token string and returns its claims.
func (m *Manager) Parse(tokenStr string) (*Claims, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenStr, claims, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return m.secret, nil
	})
	if err != nil || !token.Valid {
		return nil, errors.New("invalid token")
	}
	return claims, nil
}
