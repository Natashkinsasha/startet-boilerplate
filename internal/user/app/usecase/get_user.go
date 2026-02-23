package usecase

import (
	apperror "starter-boilerplate/internal/shared/error"
	"starter-boilerplate/internal/shared/middleware"
	"starter-boilerplate/internal/user/app/service"
	"starter-boilerplate/internal/user/domain/model"
)

type GetUserUseCase struct {
	userService service.UserService
}

func NewGetUserUseCase(us service.UserService) *GetUserUseCase {
	return &GetUserUseCase{userService: us}
}

func (uc *GetUserUseCase) Execute(ctx middleware.AuthCtx, targetID string) (*model.User, error) {
	claims := ctx.Claims()
	if claims.Role != "admin" && claims.UserID != targetID {
		return nil, apperror.ErrAccessDenied
	}

	u, err := uc.userService.FindByID(ctx, targetID)
	if err != nil {
		return nil, err
	}
	if u == nil {
		return nil, apperror.ErrNotFound
	}
	return u, nil
}
