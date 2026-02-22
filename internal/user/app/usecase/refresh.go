package usecase

import (
	"context"

	apperror "starter-boilerplate/internal/shared/error"
	"starter-boilerplate/internal/user/app/service"
	"starter-boilerplate/internal/user/domain/model"
)

type RefreshUseCase struct {
	userService  service.UserService
	tokenService service.TokenService
}

func NewRefreshUseCase(us service.UserService, ts service.TokenService) *RefreshUseCase {
	return &RefreshUseCase{
		userService:  us,
		tokenService: ts,
	}
}

func (uc *RefreshUseCase) Execute(ctx context.Context, refreshToken string) (*model.TokenPair, error) {
	claims, err := uc.tokenService.ValidateRefreshToken(refreshToken)
	if err != nil {
		return nil, err
	}

	u, err := uc.userService.FindByID(ctx, claims.UserID)
	if err != nil {
		return nil, err
	}
	if u == nil {
		return nil, apperror.ErrNotFound
	}

	return uc.tokenService.IssueTokenPair(u.ID, string(u.Role))
}
