package usecase

import (
	"context"

	"starter-boilerplate/internal/shared/errs"
	"starter-boilerplate/internal/user/app/service"
	domainevent "starter-boilerplate/internal/user/domain/event"
	"starter-boilerplate/internal/user/domain/model"
	"starter-boilerplate/pkg/outbox"
)

type LoginUseCase struct {
	userService  service.UserService
	tokenService service.TokenService
	bus          outbox.Bus
}

func NewLoginUseCase(us service.UserService, ts service.TokenService, bus outbox.Bus) *LoginUseCase {
	return &LoginUseCase{
		userService:  us,
		tokenService: ts,
		bus:          bus,
	}
}

func (uc *LoginUseCase) Execute(ctx context.Context, email, password, ip, userAgent string) (*model.TokenPair, error) {
	u, err := uc.userService.FindByEmail(ctx, email)
	if err != nil {
		return nil, err
	}
	if u == nil {
		return nil, errs.ErrInvalidCredentials
	}

	if err := uc.userService.CheckPassword(u.PasswordHash, password); err != nil {
		return nil, err
	}

	err = uc.bus.Publish(ctx, domainevent.UserLoggedInEvent{
		UserID:    u.ID,
		IP:        ip,
		UserAgent: userAgent,
	})
	if err != nil {
		return nil, err
	}

	return uc.tokenService.IssueTokenPair(u.ID, string(u.Role))
}
