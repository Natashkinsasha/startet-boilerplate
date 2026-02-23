package usecase

import (
	"starter-boilerplate/internal/shared/errs"
	"starter-boilerplate/internal/shared/middleware"
	"starter-boilerplate/internal/user/app/service"
	domainevent "starter-boilerplate/internal/user/domain/event"
	"starter-boilerplate/pkg/event"
)

type ChangePasswordUseCase struct {
	userService service.UserService
	eventBus    event.Bus
}

func NewChangePasswordUseCase(us service.UserService, eb event.Bus) *ChangePasswordUseCase {
	return &ChangePasswordUseCase{userService: us, eventBus: eb}
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

	if err := uc.userService.UpdatePassword(ctx, userID, hash); err != nil {
		return err
	}

	return uc.eventBus.Publish(ctx, domainevent.PasswordChangedEvent{
		UserID: userID,
	})
}
