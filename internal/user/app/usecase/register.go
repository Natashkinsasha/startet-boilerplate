package usecase

import (
	"context"

	"github.com/google/uuid"

	"starter-boilerplate/internal/shared/errs"
	"starter-boilerplate/internal/user/app/service"
	domainevent "starter-boilerplate/internal/user/domain/event"
	"starter-boilerplate/internal/user/domain/model"
	pkgdb "starter-boilerplate/pkg/db"
	"starter-boilerplate/pkg/outbox"
)

type RegisterUseCase struct {
	userService  service.UserService
	tokenService service.TokenService
	bus          outbox.Bus
	uow          pkgdb.UoW
}

func NewRegisterUseCase(us service.UserService, ts service.TokenService, bus outbox.Bus, uow pkgdb.UoW) *RegisterUseCase {
	return &RegisterUseCase{
		userService:  us,
		tokenService: ts,
		bus:          bus,
		uow:          uow,
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

	err = uc.uow.Do(ctx, func(ctx context.Context) error {
		if err := uc.userService.Create(ctx, user); err != nil {
			return err
		}
		return uc.bus.Publish(ctx, domainevent.UserCreatedEvent{
			UserID: user.ID,
			Email:  user.Email,
		})
	})
	if err != nil {
		return nil, err
	}

	return uc.tokenService.IssueTokenPair(user.ID, string(user.Role))
}
