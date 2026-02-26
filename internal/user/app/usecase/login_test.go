//go:build unit

package usecase

import (
	"context"
	"errors"
	"testing"

	servicemocks "starter-boilerplate/internal/user/app/service/mocks"
	"starter-boilerplate/internal/user/domain/model"
	"starter-boilerplate/pkg/outbox"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// --- bus mock ---

type mockBus struct {
	mock.Mock
}

func (m *mockBus) Publish(ctx context.Context, event outbox.Event) error {
	args := m.Called(ctx, event)
	return args.Error(0)
}

// --- tests ---

func TestLogin_Success(t *testing.T) {
	userSvc := new(servicemocks.UserService)
	tokenSvc := new(servicemocks.TokenService)
	bus := new(mockBus)
	uc := NewLoginUseCase(userSvc, tokenSvc, bus)

	user := &model.User{ID: "1", Email: "test@example.com", PasswordHash: "hash", Role: model.RoleUser}
	pair := &model.TokenPair{AccessToken: "access", RefreshToken: "refresh"}

	userSvc.On("FindByEmail", mock.Anything, "test@example.com").Return(user, nil)
	userSvc.On("CheckPassword", "hash", "password123").Return(nil)
	bus.On("Publish", mock.Anything, mock.Anything).Return(nil)
	tokenSvc.On("IssueTokenPair", "1", "user").Return(pair, nil)

	result, err := uc.Execute(context.Background(), "test@example.com", "password123", "1.2.3.4", "TestAgent/1.0")

	assert.NoError(t, err)
	assert.Equal(t, pair, result)
	userSvc.AssertExpectations(t)
	tokenSvc.AssertExpectations(t)
	bus.AssertExpectations(t)
}

func TestLogin_UserNotFound(t *testing.T) {
	userSvc := new(servicemocks.UserService)
	tokenSvc := new(servicemocks.TokenService)
	uc := NewLoginUseCase(userSvc, tokenSvc, nil)

	userSvc.On("FindByEmail", mock.Anything, "missing@example.com").Return(nil, nil)

	result, err := uc.Execute(context.Background(), "missing@example.com", "password123", "", "")

	assert.Nil(t, result)
	assert.EqualError(t, err, "invalid credentials")
	tokenSvc.AssertNotCalled(t, "IssueTokenPair")
}

func TestLogin_WrongPassword(t *testing.T) {
	userSvc := new(servicemocks.UserService)
	tokenSvc := new(servicemocks.TokenService)
	uc := NewLoginUseCase(userSvc, tokenSvc, nil)

	user := &model.User{ID: "1", Email: "test@example.com", PasswordHash: "hash", Role: model.RoleUser}

	userSvc.On("FindByEmail", mock.Anything, "test@example.com").Return(user, nil)
	userSvc.On("CheckPassword", "hash", "wrong").Return(errors.New("invalid credentials"))

	result, err := uc.Execute(context.Background(), "test@example.com", "wrong", "", "")

	assert.Nil(t, result)
	assert.EqualError(t, err, "invalid credentials")
	tokenSvc.AssertNotCalled(t, "IssueTokenPair")
}

func TestLogin_RepoError(t *testing.T) {
	userSvc := new(servicemocks.UserService)
	tokenSvc := new(servicemocks.TokenService)
	uc := NewLoginUseCase(userSvc, tokenSvc, nil)

	userSvc.On("FindByEmail", mock.Anything, "test@example.com").Return(nil, errors.New("db error"))

	result, err := uc.Execute(context.Background(), "test@example.com", "password123", "", "")

	assert.Nil(t, result)
	assert.EqualError(t, err, "db error")
	tokenSvc.AssertNotCalled(t, "IssueTokenPair")
}
