package middleware

import (
	"github.com/danielgtaylor/huma/v2"
)

func NewRoleMiddleware(api huma.API) func(huma.Context, func(huma.Context)) {
	return func(ctx huma.Context, next func(huma.Context)) {
		roles, ok := ctx.Operation().Metadata["requiredRoles"].([]string)
		if !ok || len(roles) == 0 {
			next(ctx)
			return
		}

		claims := ClaimsFromContext(ctx.Context())
		if claims == nil {
			_ = huma.WriteErr(api, ctx, 401, "unauthorized")
			return
		}

		for _, r := range roles {
			if claims.Role == r {
				next(ctx)
				return
			}
		}

		_ = huma.WriteErr(api, ctx, 403, "insufficient permissions")
	}
}
