package mocks

import (
	"starter-boilerplate/internal/user/domain/model"
	"starter-boilerplate/pkg/jwt"

	"github.com/stretchr/testify/mock"
)

type TokenService struct {
	mock.Mock
}

func (m *TokenService) IssueTokenPair(userID, role string) (*model.TokenPair, error) {
	args := m.Called(userID, role)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.TokenPair), args.Error(1)
}

func (m *TokenService) ValidateRefreshToken(token string) (*jwt.Claims, error) {
	args := m.Called(token)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*jwt.Claims), args.Error(1)
}
