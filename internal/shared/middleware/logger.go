package middleware

import (
	"log/slog"
	"net"
	"time"

	"github.com/danielgtaylor/huma/v2"
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
		attrs := []any{
			slog.String("method", ctx.Method()),
			slog.String("path", u.Path),
			slog.Int("status", status),
			slog.Duration("latency", time.Since(start)),
			slog.String("ip", ip),
			slog.String("request_id", RequestIDFromContext(ctx.Context())),
		}

		switch {
		case status >= 500:
			slog.Error("http request", attrs...)
		case status >= 400:
			slog.Warn("http request", attrs...)
		default:
			slog.Info("http request", attrs...)
		}
	}
}
