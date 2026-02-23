package service

import (
	"starter-boilerplate/internal/shared/errs"
	"starter-boilerplate/internal/user/domain/model"
	"starter-boilerplate/pkg/jwt"
)

type TokenService interface {
	IssueTokenPair(userID, role string) (*model.TokenPair, error)
	ValidateRefreshToken(token string) (*jwt.Claims, error)
}

type tokenService struct {
	jwtManager *jwt.Manager
}

func NewTokenService(jwtManager *jwt.Manager) TokenService {
	return &tokenService{jwtManager: jwtManager}
}

func (s *tokenService) IssueTokenPair(userID, role string) (*model.TokenPair, error) {
	accessToken, err := s.jwtManager.GenerateAccessToken(userID, role)
	if err != nil {
		return nil, err
	}

	refreshToken, err := s.jwtManager.GenerateRefreshToken(userID, role)
	if err != nil {
		return nil, err
	}

	return &model.TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (s *tokenService) ValidateRefreshToken(token string) (*jwt.Claims, error) {
	claims, err := s.jwtManager.ValidateRefreshToken(token)
	if err != nil {
		return nil, errs.ErrInvalidToken
	}
	return claims, nil
}
