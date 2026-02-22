//go:build unit

package service

import (
	"context"
	"errors"
	"testing"

	"starter-boilerplate/internal/user/domain/model"
	"starter-boilerplate/internal/user/domain/repository/mocks"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
)

func TestFindByEmail_Success(t *testing.T) {
	repo := new(mocks.UserRepository)
	svc := NewUserService(repo)

	expected := &model.User{ID: "1", Email: "test@example.com"}
	repo.On("FindByEmail", context.Background(), "test@example.com").Return(expected, nil)

	user, err := svc.FindByEmail(context.Background(), "test@example.com")

	require.NoError(t, err)
	assert.Equal(t, expected, user)
	repo.AssertExpectations(t)
}

func TestFindByEmail_NotFound(t *testing.T) {
	repo := new(mocks.UserRepository)
	svc := NewUserService(repo)

	repo.On("FindByEmail", context.Background(), "missing@example.com").Return(nil, nil)

	user, err := svc.FindByEmail(context.Background(), "missing@example.com")

	require.NoError(t, err)
	assert.Nil(t, user)
	repo.AssertExpectations(t)
}

func TestFindByEmail_Error(t *testing.T) {
	repo := new(mocks.UserRepository)
	svc := NewUserService(repo)

	repo.On("FindByEmail", context.Background(), "test@example.com").Return(nil, errors.New("db error"))

	user, err := svc.FindByEmail(context.Background(), "test@example.com")

	assert.Error(t, err)
	assert.Nil(t, user)
	repo.AssertExpectations(t)
}

func TestFindByID_Success(t *testing.T) {
	repo := new(mocks.UserRepository)
	svc := NewUserService(repo)

	expected := &model.User{ID: "1", Email: "test@example.com"}
	repo.On("FindByID", context.Background(), "1").Return(expected, nil)

	user, err := svc.FindByID(context.Background(), "1")

	require.NoError(t, err)
	assert.Equal(t, expected, user)
	repo.AssertExpectations(t)
}

func TestCheckPassword_Valid(t *testing.T) {
	repo := new(mocks.UserRepository)
	svc := NewUserService(repo)

	hash, err := bcrypt.GenerateFromPassword([]byte("correct-password"), bcrypt.DefaultCost)
	require.NoError(t, err)

	err = svc.CheckPassword(string(hash), "correct-password")

	assert.NoError(t, err)
}

func TestCheckPassword_Invalid(t *testing.T) {
	repo := new(mocks.UserRepository)
	svc := NewUserService(repo)

	hash, err := bcrypt.GenerateFromPassword([]byte("correct-password"), bcrypt.DefaultCost)
	require.NoError(t, err)

	err = svc.CheckPassword(string(hash), "wrong-password")

	assert.EqualError(t, err, "invalid credentials")
}
