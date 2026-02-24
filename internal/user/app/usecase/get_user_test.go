//go:build unit

package usecase

import (
	"context"
	"errors"
	"testing"

	"starter-boilerplate/internal/shared/errs"
	servicemocks "starter-boilerplate/internal/user/app/service/mocks"
	"starter-boilerplate/internal/user/domain/model"
	"starter-boilerplate/pkg/jwt"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type stubAuthCtx struct {
	context.Context
	claims *jwt.Claims
}

func (s *stubAuthCtx) Claims() *jwt.Claims { return s.claims }

func newAuthCtx(userID, role string) *stubAuthCtx {
	return &stubAuthCtx{
		Context: context.Background(),
		claims:  &jwt.Claims{UserID: userID, Role: role},
	}
}

func TestGetUser_AdminCanAccessAnyUser(t *testing.T) {
	userSvc := new(servicemocks.UserService)
	uc := NewGetUserUseCase(userSvc)

	target := &model.User{ID: "target-1", Email: "target@example.com", Role: model.RoleUser}
	userSvc.On("FindByID", mock.Anything, "target-1").Return(target, nil)

	result, err := uc.Execute(newAuthCtx("admin-1", "admin"), "target-1")

	assert.NoError(t, err)
	assert.Equal(t, target, result)
	userSvc.AssertExpectations(t)
}

func TestGetUser_UserCanAccessSelf(t *testing.T) {
	userSvc := new(servicemocks.UserService)
	uc := NewGetUserUseCase(userSvc)

	self := &model.User{ID: "user-1", Email: "self@example.com", Role: model.RoleUser}
	userSvc.On("FindByID", mock.Anything, "user-1").Return(self, nil)

	result, err := uc.Execute(newAuthCtx("user-1", "user"), "user-1")

	assert.NoError(t, err)
	assert.Equal(t, self, result)
	userSvc.AssertExpectations(t)
}

func TestGetUser_UserCannotAccessOther(t *testing.T) {
	userSvc := new(servicemocks.UserService)
	uc := NewGetUserUseCase(userSvc)

	result, err := uc.Execute(newAuthCtx("user-1", "user"), "other-1")

	assert.Nil(t, result)
	assert.ErrorIs(t, err, errs.ErrAccessDenied)
	userSvc.AssertNotCalled(t, "FindByID")
}

func TestGetUser_NotFound(t *testing.T) {
	userSvc := new(servicemocks.UserService)
	uc := NewGetUserUseCase(userSvc)

	userSvc.On("FindByID", mock.Anything, "missing-1").Return(nil, nil)

	result, err := uc.Execute(newAuthCtx("missing-1", "user"), "missing-1")

	assert.Nil(t, result)
	assert.ErrorIs(t, err, errs.ErrNotFound)
	userSvc.AssertExpectations(t)
}

func TestGetUser_RepoError(t *testing.T) {
	userSvc := new(servicemocks.UserService)
	uc := NewGetUserUseCase(userSvc)

	userSvc.On("FindByID", mock.Anything, "user-1").Return(nil, errors.New("db error"))

	result, err := uc.Execute(newAuthCtx("user-1", "user"), "user-1")

	assert.Nil(t, result)
	assert.EqualError(t, err, "db error")
	userSvc.AssertExpectations(t)
}
