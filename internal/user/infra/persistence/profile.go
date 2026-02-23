package persistence

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
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

	_, err := r.db.NewInsert().
		Model(m).
		On("CONFLICT (user_id) DO NOTHING").
		Exec(ctx)
	return err
}

func (r *profileRepository) Update(ctx context.Context, userID string, upd *model.ProfileUpdate) error {
	q := r.db.NewUpdate().
		TableExpr("user_profiles").
		Where("user_id = ?", userID).
		Set("updated_at = ?", time.Now().Unix())

	if expr, args := buildNumbersExpr(upd); expr != "" {
		q = q.Set(expr, args...)
	}
	if expr, args := buildStringsExpr(upd); expr != "" {
		q = q.Set(expr, args...)
	}

	_, err := q.Exec(ctx)
	return err
}

// buildNumbersExpr builds a SET expression for the numbers JSONB column.
// Every operation uses per-key jsonb_set so concurrent updates to different keys
// don't interfere, and increments are atomic (read-modify-write in a single expression).
func buildNumbersExpr(upd *model.ProfileUpdate) (string, []interface{}) {
	if len(upd.NumberSets) == 0 && len(upd.NumberIncr) == 0 {
		return "", nil
	}

	expr := "numbers"
	var args []interface{}

	for key, val := range upd.NumberSets {
		expr = fmt.Sprintf("jsonb_set(%s, ?::text[], to_jsonb(?::numeric))", expr)
		args = append(args, "{"+key+"}", val)
	}

	for key, delta := range upd.NumberIncr {
		expr = fmt.Sprintf(
			"jsonb_set(%s, ?::text[], to_jsonb(COALESCE((numbers->>?)::numeric, 0) + ?::numeric))",
			expr,
		)
		args = append(args, "{"+key+"}", key, delta)
	}

	return "numbers = " + expr, args
}

// buildStringsExpr builds a SET expression for the strings JSONB column.
// Uses per-key jsonb_set for safe concurrent modifications.
func buildStringsExpr(upd *model.ProfileUpdate) (string, []interface{}) {
	if len(upd.StringSets) == 0 {
		return "", nil
	}

	expr := "strings"
	var args []interface{}

	for key, val := range upd.StringSets {
		expr = fmt.Sprintf("jsonb_set(%s, ?::text[], to_jsonb(?::text))", expr)
		args = append(args, "{"+key+"}", val)
	}

	return "strings = " + expr, args
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
