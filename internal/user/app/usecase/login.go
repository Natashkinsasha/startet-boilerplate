package usecase

import (
	"context"
	"errors"

	"starter-boilerplate/internal/user/app/service"
	"starter-boilerplate/internal/user/domain/model"
)

type LoginUseCase struct {
	userService  service.UserService
	tokenService service.TokenService
}

func NewLoginUseCase(us service.UserService, ts service.TokenService) *LoginUseCase {
	return &LoginUseCase{
		userService:  us,
		tokenService: ts,
	}
}

func (uc *LoginUseCase) Execute(ctx context.Context, email, password string) (*model.TokenPair, error) {
	u, err := uc.userService.FindByEmail(ctx, email)
	if err != nil {
		return nil, err
	}
	if u == nil {
		return nil, errors.New("invalid credentials")
	}

	if err := uc.userService.CheckPassword(u.PasswordHash, password); err != nil {
		return nil, err
	}

	return uc.tokenService.IssueTokenPair(u.ID, string(u.Role))
}
