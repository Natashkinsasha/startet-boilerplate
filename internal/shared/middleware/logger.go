package middleware

import (
	"net"
	"time"

	"github.com/danielgtaylor/huma/v2"
	"go.uber.org/zap"
)

func newLoggerMiddleware() func(huma.Context, func(huma.Context)) {
	return func(ctx huma.Context, next func(huma.Context)) {
		start := time.Now()

		next(ctx)

		status := ctx.Status()
		if status == 0 {
			status = 200
		}

		ip, _, err := net.SplitHostPort(ctx.RemoteAddr())
		if err != nil {
			ip = ctx.RemoteAddr()
		}

		u := ctx.URL()
		fields := []zap.Field{
			zap.String("method", ctx.Method()),
			zap.String("path", u.Path),
			zap.Int("status", status),
			zap.Duration("latency", time.Since(start)),
			zap.String("ip", ip),
			zap.String("request_id", RequestIDFromContext(ctx.Context())),
		}

		switch {
		case status >= 500:
			zap.L().Error("http request", fields...)
		case status >= 400:
			zap.L().Warn("http request", fields...)
		default:
			zap.L().Info("http request", fields...)
		}
	}
}
