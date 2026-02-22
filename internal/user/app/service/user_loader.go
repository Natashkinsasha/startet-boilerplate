package service

import (
	"context"
	"sync"

	"starter-boilerplate/internal/user/domain/repository"
)

type AuthUser struct {
	ID    string
	Email string
	Role  string
}

type UserLoaderCreator struct {
	repo repository.UserRepository
}

func NewUserLoaderCreator(repo repository.UserRepository) *UserLoaderCreator {
	return &UserLoaderCreator{repo: repo}
}

func (c *UserLoaderCreator) Create(userID string) func(ctx context.Context) (*AuthUser, error) {
	var (
		once   sync.Once
		cached *AuthUser
		err    error
	)

	return func(ctx context.Context) (*AuthUser, error) {
		once.Do(func() {
			u, findErr := c.repo.FindByID(ctx, userID)
			if findErr != nil {
				err = findErr
				return
			}
			cached = &AuthUser{
				ID:    u.ID,
				Email: u.Email,
				Role:  string(u.Role),
			}
		})
		return cached, err
	}
}
