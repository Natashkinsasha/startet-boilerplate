package middleware

import (
	"log/slog"
	"net/http"
	"time"

	"starter-boilerplate/internal/user/app/service"
	"starter-boilerplate/pkg/jwt"

	"github.com/danielgtaylor/huma/v2"
)

type Init struct{}

func Setup(srv *http.Server, api huma.API, jwtManager *jwt.Manager, userLoaderCreator *service.UserLoaderCreator) Init {
	// Huma-level middleware (order: outermost first)
	api.UseMiddleware(NewRequestIDMiddleware())
	api.UseMiddleware(newLoggerMiddleware())
	api.UseMiddleware(NewLimiterMiddleware(100, time.Minute))
	api.UseMiddleware(NewAuthMiddleware(api, jwtManager))
	api.UseMiddleware(NewRoleMiddleware(api))
	api.UseMiddleware(NewUserLoaderMiddleware(userLoaderCreator))

	// HTTP-level middleware
	srv.Handler = WithCORS(WithRecover(srv.Handler))

	slog.Info("middlewares installed")

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
				slog.Error("panic recovered",
					slog.Any("error", err),
					slog.String("method", r.Method),
					slog.String("path", r.URL.Path),
				)
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusInternalServerError)
				_, _ = w.Write([]byte(`{"title":"Internal Server Error","status":500}`))
			}
		}()
		next.ServeHTTP(w, r)
	})
}
