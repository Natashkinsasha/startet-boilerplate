//go:build unit

package usecase

import (
	"context"
	"errors"
	"testing"

	servicemocks "starter-boilerplate/internal/user/app/service/mocks"
	"starter-boilerplate/internal/user/domain/model"
	"starter-boilerplate/pkg/jwt"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestRefresh_Success(t *testing.T) {
	userSvc := new(servicemocks.UserService)
	tokenSvc := new(servicemocks.TokenService)
	uc := NewRefreshUseCase(userSvc, tokenSvc)

	claims := &jwt.Claims{UserID: "1", Role: "user"}
	user := &model.User{ID: "1", Email: "test@example.com", Role: model.RoleUser}
	pair := &model.TokenPair{AccessToken: "new-access", RefreshToken: "new-refresh"}

	tokenSvc.On("ValidateRefreshToken", "valid-token").Return(claims, nil)
	userSvc.On("FindByID", mock.Anything, "1").Return(user, nil)
	tokenSvc.On("IssueTokenPair", "1", "user").Return(pair, nil)

	result, err := uc.Execute(context.Background(), "valid-token")

	assert.NoError(t, err)
	assert.Equal(t, pair, result)
	userSvc.AssertExpectations(t)
	tokenSvc.AssertExpectations(t)
}

func TestRefresh_InvalidToken(t *testing.T) {
	userSvc := new(servicemocks.UserService)
	tokenSvc := new(servicemocks.TokenService)
	uc := NewRefreshUseCase(userSvc, tokenSvc)

	tokenSvc.On("ValidateRefreshToken", "garbage").Return(nil, errors.New("invalid refresh token"))

	result, err := uc.Execute(context.Background(), "garbage")

	assert.Nil(t, result)
	assert.EqualError(t, err, "invalid refresh token")
	userSvc.AssertNotCalled(t, "FindByID")
}

func TestRefresh_UserNotFound(t *testing.T) {
	userSvc := new(servicemocks.UserService)
	tokenSvc := new(servicemocks.TokenService)
	uc := NewRefreshUseCase(userSvc, tokenSvc)

	claims := &jwt.Claims{UserID: "999", Role: "user"}

	tokenSvc.On("ValidateRefreshToken", "valid-token").Return(claims, nil)
	userSvc.On("FindByID", mock.Anything, "999").Return(nil, nil)

	result, err := uc.Execute(context.Background(), "valid-token")

	assert.Nil(t, result)
	assert.EqualError(t, err, "not found")
	tokenSvc.AssertNotCalled(t, "IssueTokenPair")
}
