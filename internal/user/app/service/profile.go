package service

import (
	"context"

	domainevent "starter-boilerplate/internal/user/domain/event"
	"starter-boilerplate/internal/user/domain/model"
	"starter-boilerplate/internal/user/domain/repository"
	pkgamqp "starter-boilerplate/pkg/amqp"
)

type ProfileService struct {
	profileRepo repository.ProfileRepository
}

func NewProfileService(pr repository.ProfileRepository) *ProfileService {
	return &ProfileService{profileRepo: pr}
}

func (s *ProfileService) OnUserCreated(ctx context.Context, evt domainevent.UserCreatedEvent, _ pkgamqp.DeliveryMeta) error {
	return s.profileRepo.Upsert(ctx, &model.Profile{
		UserID:  evt.UserID,
		Numbers: map[string]float64{},
		Strings: map[string]string{},
	})
}

func (s *ProfileService) OnPasswordChanged(ctx context.Context, evt domainevent.PasswordChangedEvent, _ pkgamqp.DeliveryMeta) error {
	upd := model.NewProfileUpdate().
		IncrNumber("password_changes", 1)

	return s.profileRepo.Update(ctx, evt.UserID, upd)
}
