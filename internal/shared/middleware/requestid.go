package middleware

import (
	"context"

	"github.com/danielgtaylor/huma/v2"
	"github.com/google/uuid"
)

type requestIDKey struct{}

func NewRequestIDMiddleware() func(huma.Context, func(huma.Context)) {
	return func(ctx huma.Context, next func(huma.Context)) {
		id := ctx.Header("X-Request-Id")
		if id == "" {
			id = uuid.New().String()
		}
		ctx.SetHeader("X-Request-Id", id)
		ctx = huma.WithValue(ctx, requestIDKey{}, id)
		next(ctx)
	}
}

func RequestIDFromContext(ctx context.Context) string {
	id, _ := ctx.Value(requestIDKey{}).(string)
	return id
}
