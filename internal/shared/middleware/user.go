package middleware

import (
	"context"
	"sync"

	"starter-boilerplate/internal/user/app/service"

	"github.com/danielgtaylor/huma/v2"
)

type userLoaderKey struct{}

func NewUserLoaderMiddleware(creator *service.UserLoaderCreator) func(huma.Context, func(huma.Context)) {
	return func(ctx huma.Context, next func(huma.Context)) {
		claims, ok := AuthFromCtx(ctx.Context())
		if ok {
			loader := creator.Create(claims.UserID)
			ctx = huma.WithValue(ctx, userLoaderKey{}, loader)
		}
		next(ctx)
	}
}

func UserLoaderFromContext(ctx context.Context) (func(context.Context) (*service.AuthUser, error), bool) {
	loader, ok := ctx.Value(userLoaderKey{}).(func(context.Context) (*service.AuthUser, error))
	return loader, ok && loader != nil
}

type UserCtx interface {
	AuthCtx
	User() (*service.AuthUser, error)
}

type userCtx struct {
	AuthCtx
	userOnce sync.Once
	user     *service.AuthUser
	userErr  error
}

func NewUserCtx(ctx context.Context) UserCtx {
	return &userCtx{AuthCtx: NewAuthCtx(ctx)}
}

func (a *userCtx) User() (*service.AuthUser, error) {
	a.userOnce.Do(func() {
		loader, ok := UserLoaderFromContext(a)
		if !ok {
			panic("UserCtx.User: no user loader in context")
		}
		a.user, a.userErr = loader(a.AuthCtx)
	})
	return a.user, a.userErr
}
