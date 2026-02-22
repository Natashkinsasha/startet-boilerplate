package middleware

import (
	"context"

	"starter-boilerplate/internal/user/app/service"

	"github.com/danielgtaylor/huma/v2"
)

type userLoaderKey struct{}

func NewUserLoaderMiddleware(creator *service.UserLoaderCreator) func(huma.Context, func(huma.Context)) {
	return func(ctx huma.Context, next func(huma.Context)) {
		claims, ok := ClaimsFromContextSafe(ctx.Context())
		if ok {
			loader := creator.Create(claims.UserID)
			ctx = huma.WithValue(ctx, userLoaderKey{}, loader)
		}
		next(ctx)
	}
}

func UserLoaderFromContext(ctx context.Context) func(context.Context) (*service.AuthUser, error) {
	loader, _ := ctx.Value(userLoaderKey{}).(func(context.Context) (*service.AuthUser, error))
	return loader
}
