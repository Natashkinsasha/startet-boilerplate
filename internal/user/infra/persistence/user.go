package persistence

import (
	"context"
	"database/sql"
	"errors"

	"starter-boilerplate/internal/user/domain/model"
	"starter-boilerplate/internal/user/domain/repository"
	pkgdb "starter-boilerplate/pkg/db"

	"github.com/uptrace/bun"
)

type userModel struct {
	bun.BaseModel `bun:"table:users"`

	ID           string `bun:"id,pk"`
	Email        string `bun:"email,unique,notnull"`
	PasswordHash string `bun:"password_hash,notnull"`
	Role         string `bun:"role,notnull,default:'user'"`
	CreatedAt    int64  `bun:"created_at,notnull"`
	UpdatedAt    int64  `bun:"updated_at,notnull"`
}

type userRepository struct {
	db *bun.DB
}

func NewUserRepository(db *bun.DB) repository.UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) FindByID(ctx context.Context, id string) (*model.User, error) {
	var m userModel
	err := pkgdb.Conn(ctx, r.db).NewSelect().Model(&m).Where("id = ?", id).Scan(ctx)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return toEntity(&m), nil
}

func (r *userRepository) FindByEmail(ctx context.Context, email string) (*model.User, error) {
	var m userModel
	err := pkgdb.Conn(ctx, r.db).NewSelect().Model(&m).Where("email = ?", email).Scan(ctx)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return toEntity(&m), nil
}

func (r *userRepository) Create(ctx context.Context, u *model.User) error {
	m := fromEntity(u)
	_, err := pkgdb.Conn(ctx, r.db).NewInsert().Model(m).Exec(ctx)
	return err
}

func (r *userRepository) Update(ctx context.Context, u *model.User) error {
	m := fromEntity(u)
	_, err := pkgdb.Conn(ctx, r.db).NewUpdate().Model(m).WherePK().Exec(ctx)
	return err
}

func (r *userRepository) UpdatePassword(ctx context.Context, id, hash string) error {
	_, err := pkgdb.Conn(ctx, r.db).NewUpdate().
		Model((*userModel)(nil)).
		Set("password_hash = ?", hash).
		Where("id = ?", id).
		Exec(ctx)
	return err
}

func toEntity(m *userModel) *model.User {
	return &model.User{
		ID:           m.ID,
		Email:        m.Email,
		PasswordHash: m.PasswordHash,
		Role:         model.Role(m.Role),
	}
}

func fromEntity(u *model.User) *userModel {
	return &userModel{
		ID:           u.ID,
		Email:        u.Email,
		PasswordHash: u.PasswordHash,
		Role:         string(u.Role),
	}
}
