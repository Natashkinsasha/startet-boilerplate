package repository

import (
	"context"

	"starter-boilerplate/internal/user/domain/model"
)

type ProfileRepository interface {
	FindByUserID(ctx context.Context, userID string) (*model.Profile, error)
	Upsert(ctx context.Context, profile *model.Profile) error
	Update(ctx context.Context, userID string, upd *model.ProfileUpdate) error
}
