package usecase

import (
	"context"

	"starter-boilerplate/internal/shared/errs"
	"starter-boilerplate/internal/shared/middleware"
	"starter-boilerplate/internal/user/app/service"
	domainevent "starter-boilerplate/internal/user/domain/event"
	pkgdb "starter-boilerplate/pkg/db"
	"starter-boilerplate/pkg/outbox"
)

type ChangePasswordUseCase struct {
	userService service.UserService
	bus         outbox.Bus
	uow         pkgdb.UoW
}

func NewChangePasswordUseCase(us service.UserService, bus outbox.Bus, uow pkgdb.UoW) *ChangePasswordUseCase {
	return &ChangePasswordUseCase{userService: us, bus: bus, uow: uow}
}

func (uc *ChangePasswordUseCase) Execute(ctx middleware.AuthCtx, oldPassword, newPassword string) error {
	userID := ctx.Claims().UserID

	user, err := uc.userService.FindByID(ctx, userID)
	if err != nil {
		return err
	}
	if user == nil {
		return errs.ErrNotFound
	}

	if err := uc.userService.CheckPassword(user.PasswordHash, oldPassword); err != nil {
		return err
	}

	hash, err := uc.userService.HashPassword(newPassword)
	if err != nil {
		return err
	}

	return uc.uow.Do(ctx, func(ctx context.Context) error {
		if err := uc.userService.UpdatePassword(ctx, userID, hash); err != nil {
			return err
		}
		return uc.bus.Publish(ctx, domainevent.PasswordChangedEvent{
			UserID: userID,
		})
	})
}
