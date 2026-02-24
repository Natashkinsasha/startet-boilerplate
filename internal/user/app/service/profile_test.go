//go:build unit

package service

import (
	"context"
	"errors"
	"testing"

	pkgamqp "starter-boilerplate/pkg/amqp"

	domainevent "starter-boilerplate/internal/user/domain/event"
	"starter-boilerplate/internal/user/domain/model"
	repomocks "starter-boilerplate/internal/user/domain/repository/mocks"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestProfileService_OnUserCreated(t *testing.T) {
	repo := new(repomocks.ProfileRepository)
	svc := NewProfileService(repo)

	repo.On("Upsert", mock.Anything, &model.Profile{
		UserID:  "user-1",
		Numbers: map[string]float64{},
		Strings: map[string]string{},
	}).Return(nil)

	err := svc.OnUserCreated(context.Background(), domainevent.UserCreatedEvent{
		UserID: "user-1",
		Email:  "test@example.com",
	}, pkgamqp.DeliveryMeta{})

	assert.NoError(t, err)
	repo.AssertExpectations(t)
}

func TestProfileService_OnUserCreated_RepoError(t *testing.T) {
	repo := new(repomocks.ProfileRepository)
	svc := NewProfileService(repo)

	repo.On("Upsert", mock.Anything, mock.Anything).Return(errors.New("db error"))

	err := svc.OnUserCreated(context.Background(), domainevent.UserCreatedEvent{
		UserID: "user-1",
		Email:  "test@example.com",
	}, pkgamqp.DeliveryMeta{})

	assert.EqualError(t, err, "db error")
	repo.AssertExpectations(t)
}

func TestProfileService_OnPasswordChanged(t *testing.T) {
	repo := new(repomocks.ProfileRepository)
	svc := NewProfileService(repo)

	expectedUpd := model.NewProfileUpdate().IncrNumber("password_changes", 1)

	repo.On("Update", mock.Anything, "user-1", expectedUpd).Return(nil)

	err := svc.OnPasswordChanged(context.Background(), domainevent.PasswordChangedEvent{
		UserID: "user-1",
	}, pkgamqp.DeliveryMeta{})

	assert.NoError(t, err)
	repo.AssertExpectations(t)
}

func TestProfileService_OnPasswordChanged_RepoError(t *testing.T) {
	repo := new(repomocks.ProfileRepository)
	svc := NewProfileService(repo)

	repo.On("Update", mock.Anything, "user-1", mock.Anything).Return(errors.New("update failed"))

	err := svc.OnPasswordChanged(context.Background(), domainevent.PasswordChangedEvent{
		UserID: "user-1",
	}, pkgamqp.DeliveryMeta{})

	assert.EqualError(t, err, "update failed")
	repo.AssertExpectations(t)
}
