//go:build unit

package usecase

import (
	"context"
	"errors"
	"testing"

	servicemocks "starter-boilerplate/internal/user/app/service/mocks"
	"starter-boilerplate/internal/user/domain/model"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestLogin_Success(t *testing.T) {
	userSvc := new(servicemocks.UserService)
	tokenSvc := new(servicemocks.TokenService)
	uc := NewLoginUseCase(userSvc, tokenSvc)

	user := &model.User{ID: "1", Email: "test@example.com", PasswordHash: "hash", Role: model.RoleUser}
	pair := &model.TokenPair{AccessToken: "access", RefreshToken: "refresh"}

	userSvc.On("FindByEmail", mock.Anything, "test@example.com").Return(user, nil)
	userSvc.On("CheckPassword", "hash", "password123").Return(nil)
	tokenSvc.On("IssueTokenPair", "1", "user").Return(pair, nil)

	result, err := uc.Execute(context.Background(), "test@example.com", "password123")

	assert.NoError(t, err)
	assert.Equal(t, pair, result)
	userSvc.AssertExpectations(t)
	tokenSvc.AssertExpectations(t)
}

func TestLogin_UserNotFound(t *testing.T) {
	userSvc := new(servicemocks.UserService)
	tokenSvc := new(servicemocks.TokenService)
	uc := NewLoginUseCase(userSvc, tokenSvc)

	userSvc.On("FindByEmail", mock.Anything, "missing@example.com").Return(nil, nil)

	result, err := uc.Execute(context.Background(), "missing@example.com", "password123")

	assert.Nil(t, result)
	assert.EqualError(t, err, "invalid credentials")
	tokenSvc.AssertNotCalled(t, "IssueTokenPair")
}

func TestLogin_WrongPassword(t *testing.T) {
	userSvc := new(servicemocks.UserService)
	tokenSvc := new(servicemocks.TokenService)
	uc := NewLoginUseCase(userSvc, tokenSvc)

	user := &model.User{ID: "1", Email: "test@example.com", PasswordHash: "hash", Role: model.RoleUser}

	userSvc.On("FindByEmail", mock.Anything, "test@example.com").Return(user, nil)
	userSvc.On("CheckPassword", "hash", "wrong").Return(errors.New("invalid credentials"))

	result, err := uc.Execute(context.Background(), "test@example.com", "wrong")

	assert.Nil(t, result)
	assert.EqualError(t, err, "invalid credentials")
	tokenSvc.AssertNotCalled(t, "IssueTokenPair")
}

func TestLogin_RepoError(t *testing.T) {
	userSvc := new(servicemocks.UserService)
	tokenSvc := new(servicemocks.TokenService)
	uc := NewLoginUseCase(userSvc, tokenSvc)

	userSvc.On("FindByEmail", mock.Anything, "test@example.com").Return(nil, errors.New("db error"))

	result, err := uc.Execute(context.Background(), "test@example.com", "password123")

	assert.Nil(t, result)
	assert.EqualError(t, err, "db error")
	tokenSvc.AssertNotCalled(t, "IssueTokenPair")
}
