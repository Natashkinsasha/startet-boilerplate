package middleware

import (
	"context"
	"strings"
	"sync"

	"starter-boilerplate/pkg/jwt"

	"github.com/danielgtaylor/huma/v2"
)

var _ AuthCtx = (*authCtx)(nil)

type claimsContextKey struct{}

func NewAuthMiddleware(api huma.API, jwtManager *jwt.Manager) func(huma.Context, func(huma.Context)) {
	return func(ctx huma.Context, next func(huma.Context)) {
		if !requiresBearerAuth(ctx.Operation()) {
			next(ctx)
			return
		}

		header := ctx.Header("Authorization")
		if header == "" {
			_ = huma.WriteErr(api, ctx, 401, "missing authorization header")
			return
		}

		token, found := strings.CutPrefix(header, "Bearer ")
		if !found {
			_ = huma.WriteErr(api, ctx, 401, "invalid authorization header format")
			return
		}

		claims, err := jwtManager.ValidateAccessToken(token)
		if err != nil {
			_ = huma.WriteErr(api, ctx, 401, "invalid or expired token")
			return
		}

		ctx = huma.WithValue(ctx, claimsContextKey{}, claims)
		next(ctx)
	}
}

func requiresBearerAuth(op *huma.Operation) bool {
	for _, sec := range op.Security {
		if _, ok := sec["bearerAuth"]; ok {
			return true
		}
	}
	return false
}

type AuthCtx interface {
	context.Context
	Claims() *jwt.Claims
}

type authCtx struct {
	context.Context
	claimsOnce sync.Once
	claims     *jwt.Claims
}

func NewAuthCtx(ctx context.Context) AuthCtx {
	return &authCtx{Context: ctx}
}

func AuthFromCtx(ctx context.Context) (*jwt.Claims, bool) {
	claims, ok := ctx.Value(claimsContextKey{}).(*jwt.Claims)
	return claims, ok && claims != nil
}

func (c *authCtx) Claims() *jwt.Claims {
	c.claimsOnce.Do(func() {
		claims, ok := AuthFromCtx(c.Context)
		if !ok || claims == nil {
			panic("AuthCtx.Claims: no claims in context — is bearerAuth set on this operation?")
		}
		c.claims = claims
	})
	return c.claims
}
