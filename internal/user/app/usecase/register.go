package usecase

import (
	"context"

	"github.com/google/uuid"

	"starter-boilerplate/internal/shared/errs"
	"starter-boilerplate/internal/user/app/service"
	"starter-boilerplate/internal/user/domain"
	"starter-boilerplate/internal/user/domain/model"
	"starter-boilerplate/pkg/event"
)

type RegisterUseCase struct {
	userService  service.UserService
	tokenService service.TokenService
	eventBus     event.Bus
}

func NewRegisterUseCase(us service.UserService, ts service.TokenService, eb event.Bus) *RegisterUseCase {
	return &RegisterUseCase{
		userService:  us,
		tokenService: ts,
		eventBus:     eb,
	}
}

func (uc *RegisterUseCase) Execute(ctx context.Context, email, password string) (*model.TokenPair, error) {
	existing, err := uc.userService.FindByEmail(ctx, email)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return nil, errs.ErrEmailAlreadyExists
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

	if err := uc.eventBus.Publish(ctx, domain.UserCreatedEvent{
		UserID: user.ID,
		Email:  user.Email,
	}); err != nil {
		return nil, err
	}

	return uc.tokenService.IssueTokenPair(user.ID, string(user.Role))
}
