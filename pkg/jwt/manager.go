package jwt

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type tokenType string

const (
	accessToken  tokenType = "access"
	refreshToken tokenType = "refresh"
)

type Claims struct {
	UserID    string    `json:"user_id"`
	Role      string    `json:"role"`
	TokenType tokenType `json:"token_type"`
	jwt.RegisteredClaims
}

type Config struct {
	AccessSecret  string
	RefreshSecret string
	AccessTTL     time.Duration
	RefreshTTL    time.Duration
}

type Manager struct {
	cfg Config
}

func NewManager(cfg Config) *Manager {
	return &Manager{cfg: cfg}
}

func (m *Manager) GenerateAccessToken(userID, role string) (string, error) {
	return m.generate(userID, role, accessToken, m.cfg.AccessSecret, m.cfg.AccessTTL)
}

func (m *Manager) GenerateRefreshToken(userID, role string) (string, error) {
	return m.generate(userID, role, refreshToken, m.cfg.RefreshSecret, m.cfg.RefreshTTL)
}

func (m *Manager) ValidateAccessToken(tokenStr string) (*Claims, error) {
	return m.validate(tokenStr, m.cfg.AccessSecret, accessToken)
}

func (m *Manager) ValidateRefreshToken(tokenStr string) (*Claims, error) {
	return m.validate(tokenStr, m.cfg.RefreshSecret, refreshToken)
}

func (m *Manager) generate(userID, role string, tt tokenType, secret string, ttl time.Duration) (string, error) {
	claims := &Claims{
		UserID:    userID,
		Role:      role,
		TokenType: tt,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(ttl)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

func (m *Manager) validate(tokenStr, secret string, expected tokenType) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &Claims{}, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(secret), nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token")
	}

	if claims.TokenType != expected {
		return nil, errors.New("invalid token type")
	}

	return claims, nil
}
