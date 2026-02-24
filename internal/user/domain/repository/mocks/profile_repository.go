package mocks

import (
	"context"

	"starter-boilerplate/internal/user/domain/model"

	"github.com/stretchr/testify/mock"
)

type ProfileRepository struct {
	mock.Mock
}

func (m *ProfileRepository) FindByUserID(ctx context.Context, userID string) (*model.Profile, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Profile), args.Error(1)
}

func (m *ProfileRepository) Upsert(ctx context.Context, profile *model.Profile) error {
	return m.Called(ctx, profile).Error(0)
}

func (m *ProfileRepository) Update(ctx context.Context, userID string, upd *model.ProfileUpdate) error {
	return m.Called(ctx, userID, upd).Error(0)
}
