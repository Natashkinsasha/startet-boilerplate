package persistence

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"starter-boilerplate/internal/user/domain/model"
	"starter-boilerplate/internal/user/domain/repository"

	"github.com/uptrace/bun"
)

type profileModel struct {
	bun.BaseModel `bun:"table:user_profiles"`

	UserID    string             `bun:"user_id,pk"`
	Numbers   map[string]float64 `bun:"numbers,type:jsonb,notnull,default:'{}'"`
	Strings   map[string]string  `bun:"strings,type:jsonb,notnull,default:'{}'"`
	CreatedAt int64              `bun:"created_at,notnull"`
	UpdatedAt int64              `bun:"updated_at,notnull"`
}

type profileRepository struct {
	db *bun.DB
}

func NewProfileRepository(db *bun.DB) repository.ProfileRepository {
	return &profileRepository{db: db}
}

func (r *profileRepository) FindByUserID(ctx context.Context, userID string) (*model.Profile, error) {
	var m profileModel
	err := r.db.NewSelect().Model(&m).Where("user_id = ?", userID).Scan(ctx)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return toProfileEntity(&m), nil
}

func (r *profileRepository) Upsert(ctx context.Context, profile *model.Profile) error {
	m := fromProfileEntity(profile)
	now := time.Now().Unix()
	m.UpdatedAt = now

	_, err := r.db.NewInsert().
		Model(m).
		On("CONFLICT (user_id) DO UPDATE").
		Set("numbers = EXCLUDED.numbers").
		Set("strings = EXCLUDED.strings").
		Set("updated_at = EXCLUDED.updated_at").
		Exec(ctx)
	return err
}

func toProfileEntity(m *profileModel) *model.Profile {
	return &model.Profile{
		UserID:  m.UserID,
		Numbers: m.Numbers,
		Strings: m.Strings,
	}
}

func fromProfileEntity(p *model.Profile) *profileModel {
	now := time.Now().Unix()
	return &profileModel{
		UserID:    p.UserID,
		Numbers:   p.Numbers,
		Strings:   p.Strings,
		CreatedAt: now,
		UpdatedAt: now,
	}
}
