package usecase

import (
	"context"

	"github.com/google/uuid"

	apperror "starter-boilerplate/internal/shared/error"
	"starter-boilerplate/internal/user/app/service"
	"starter-boilerplate/internal/user/domain/model"
)

type RegisterUseCase struct {
	userService  service.UserService
	tokenService service.TokenService
}

func NewRegisterUseCase(us service.UserService, ts service.TokenService) *RegisterUseCase {
	return &RegisterUseCase{
		userService:  us,
		tokenService: ts,
	}
}

func (uc *RegisterUseCase) Execute(ctx context.Context, email, password string) (*model.TokenPair, error) {
	existing, err := uc.userService.FindByEmail(ctx, email)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return nil, apperror.ErrEmailAlreadyExists
	}

	hash, err := uc.userService.HashPassword(password)
	if err != nil {
		return nil, err
	}

	user := &model.User{
		ID:           uuid.New().String(),
		Email:        email,
		PasswordHash: hash,
		Role:         model.RoleUser,
	}

	if err := uc.userService.Create(ctx, user); err != nil {
		return nil, err
	}

	return uc.tokenService.IssueTokenPair(user.ID, string(user.Role))
}
