package middleware

import (
	"context"
	"strings"

	"starter-boilerplate/pkg/jwt"

	"github.com/danielgtaylor/huma/v2"
)

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

func ClaimsFromContext(ctx context.Context) *jwt.Claims {
	claims, ok := ctx.Value(claimsContextKey{}).(*jwt.Claims)
	if !ok || claims == nil {
		panic("ClaimsFromContext: no claims in context — is bearerAuth set on this operation?")
	}
	return claims
}
