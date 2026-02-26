package service

import (
	"context"

	"starter-boilerplate/internal/shared/errs"
	"starter-boilerplate/internal/user/domain/model"
	"starter-boilerplate/internal/user/domain/repository"

	"golang.org/x/crypto/bcrypt"
)

type UserService interface {
	FindByEmail(ctx context.Context, email string) (*model.User, error)
	FindByID(ctx context.Context, id string) (*model.User, error)
	CheckPassword(passwordHash, password string) error
	Create(ctx context.Context, user *model.User) error
	HashPassword(password string) (string, error)
	UpdatePassword(ctx context.Context, id, hash string) error
}

type userService struct {
	userRepo repository.UserRepository
}

func NewUserService(userRepo repository.UserRepository) UserService {
	return &userService{userRepo: userRepo}
}

func (s *userService) FindByEmail(ctx context.Context, email string) (*model.User, error) {
	return s.userRepo.FindByEmail(ctx, email)
}

func (s *userService) FindByID(ctx context.Context, id string) (*model.User, error) {
	return s.userRepo.FindByID(ctx, id)
}

func (s *userService) CheckPassword(passwordHash, password string) error {
	if err := bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(password)); err != nil {
		return errs.ErrInvalidCredentials
	}
	return nil
}

func (s *userService) HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

func (s *userService) Create(ctx context.Context, user *model.User) error {
	return s.userRepo.Create(ctx, user)
}

func (s *userService) UpdatePassword(ctx context.Context, id, hash string) error {
	return s.userRepo.UpdatePassword(ctx, id, hash)
}
