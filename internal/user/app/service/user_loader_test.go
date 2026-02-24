//go:build unit

package service

import (
	"context"
	"errors"
	"testing"

	"starter-boilerplate/internal/user/domain/model"
	repomocks "starter-boilerplate/internal/user/domain/repository/mocks"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestUserLoader_Success(t *testing.T) {
	repo := new(repomocks.UserRepository)
	creator := NewUserLoaderCreator(repo)

	user := &model.User{ID: "user-1", Email: "test@example.com", Role: model.RoleAdmin}
	repo.On("FindByID", mock.Anything, "user-1").Return(user, nil).Once()

	loader := creator.Create("user-1")

	result, err := loader(context.Background())
	require.NoError(t, err)
	assert.Equal(t, "user-1", result.ID)
	assert.Equal(t, "test@example.com", result.Email)
	assert.Equal(t, "admin", result.Role)
	repo.AssertExpectations(t)
}

func TestUserLoader_CachesResult(t *testing.T) {
	repo := new(repomocks.UserRepository)
	creator := NewUserLoaderCreator(repo)

	user := &model.User{ID: "user-1", Email: "test@example.com", Role: model.RoleUser}
	repo.On("FindByID", mock.Anything, "user-1").Return(user, nil).Once()

	loader := creator.Create("user-1")

	r1, err1 := loader(context.Background())
	r2, err2 := loader(context.Background())

	require.NoError(t, err1)
	require.NoError(t, err2)
	assert.Same(t, r1, r2)
	repo.AssertNumberOfCalls(t, "FindByID", 1)
}

func TestUserLoader_RepoError(t *testing.T) {
	repo := new(repomocks.UserRepository)
	creator := NewUserLoaderCreator(repo)

	repo.On("FindByID", mock.Anything, "user-1").Return(nil, errors.New("db error")).Once()

	loader := creator.Create("user-1")

	result, err := loader(context.Background())
	assert.Nil(t, result)
	assert.EqualError(t, err, "db error")
}

func TestUserLoader_CachesError(t *testing.T) {
	repo := new(repomocks.UserRepository)
	creator := NewUserLoaderCreator(repo)

	repo.On("FindByID", mock.Anything, "user-1").Return(nil, errors.New("db error")).Once()

	loader := creator.Create("user-1")

	_, err1 := loader(context.Background())
	_, err2 := loader(context.Background())

	assert.EqualError(t, err1, "db error")
	assert.EqualError(t, err2, "db error")
	repo.AssertNumberOfCalls(t, "FindByID", 1)
}
