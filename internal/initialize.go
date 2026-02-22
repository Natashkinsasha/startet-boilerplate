//go:build wireinject

package internal

import (
	"context"
	"net/http"

	"starter-boilerplate/internal/shared/app"
	"starter-boilerplate/internal/shared/config"
	"starter-boilerplate/internal/shared/db"
	sharedgrpc "starter-boilerplate/internal/shared/grpc"
	"starter-boilerplate/internal/shared/huma"
	sharedjwt "starter-boilerplate/internal/shared/jwt"
	"starter-boilerplate/internal/shared/logger"
	"starter-boilerplate/internal/shared/middleware"
	"starter-boilerplate/internal/shared/redis"
	"starter-boilerplate/internal/shared/server"
	"starter-boilerplate/internal/user"

	gohuma "github.com/danielgtaylor/huma/v2"
	"github.com/google/wire"
	goredis "github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	gogrpc "google.golang.org/grpc"
)

func newApp(httpSrv *http.Server, cfg *config.Config, _ user.Module, _ *zap.Logger, _ *goredis.Client, grpcSrv *gogrpc.Server, api gohuma.API) *app.App {
	return app.New(httpSrv, cfg, grpcSrv, api)
}

func InitializeApp(ctx context.Context) *app.App {
	wire.Build(
		config.SetupConfig,
		logger.SetupLogger,
		wire.FieldsOf(new(*config.Config), "App", "Logger", "DB", "JWT", "Redis", "GRPC"),
		db.Setup,
		redis.Setup,
		server.SetupMux,
		server.SetupHTTPServer,
		huma.Setup,
		sharedgrpc.Setup,
		sharedjwt.NewJWTManager,
		middleware.Setup,

		user.InitializeUserModule,

		newApp,
	)
	return nil
}
