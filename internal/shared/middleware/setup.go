package middleware

import (
	"net/http"
	"time"

	"starter-boilerplate/pkg/jwt"

	"github.com/danielgtaylor/huma/v2"
	"go.uber.org/zap"
)

type Init struct{}

func Setup(srv *http.Server, api huma.API, jwtManager *jwt.Manager) Init {
	// Huma-level middleware (order: outermost first)
	api.UseMiddleware(NewRequestIDMiddleware())
	api.UseMiddleware(newLoggerMiddleware())
	api.UseMiddleware(NewLimiterMiddleware(100, time.Minute))
	api.UseMiddleware(NewAuthMiddleware(api, jwtManager))
	api.UseMiddleware(NewRoleMiddleware(api))

	// HTTP-level middleware
	srv.Handler = WithCORS(WithRecover(srv.Handler))

	zap.L().Info("middlewares installed")

	return Init{}
}

func WithCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, PATCH, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func WithRecover(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				zap.L().Error("panic recovered",
					zap.Any("error", err),
					zap.String("method", r.Method),
					zap.String("path", r.URL.Path),
				)
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusInternalServerError)
				_, _ = w.Write([]byte(`{"title":"Internal Server Error","status":500}`))
			}
		}()
		next.ServeHTTP(w, r)
	})
}
