//go:build wireinject

package internal

import (
	"context"
	"log/slog"
	"net/http"

	"starter-boilerplate/internal/shared/app"
	"starter-boilerplate/internal/shared/config"
	"starter-boilerplate/internal/shared/huma"
	sharedjwt "starter-boilerplate/internal/shared/jwt"
	"starter-boilerplate/internal/shared/logger"
	"starter-boilerplate/internal/shared/middleware"
	"starter-boilerplate/internal/shared/server"
	"starter-boilerplate/internal/user"
	"starter-boilerplate/internal/user/app/service"
	"starter-boilerplate/internal/user/infra/persistence"
	"starter-boilerplate/internal/user/transport/consumer"
	pkgamqp "starter-boilerplate/pkg/amqp"
	"starter-boilerplate/pkg/db"
	"starter-boilerplate/pkg/event"
	pkggrpc "starter-boilerplate/pkg/grpc"
	"starter-boilerplate/pkg/redis"

	gohuma "github.com/danielgtaylor/huma/v2"
	"github.com/google/wire"
	amqp091 "github.com/rabbitmq/amqp091-go"
	goredis "github.com/redis/go-redis/v9"
	gogrpc "google.golang.org/grpc"
)

func newApp(httpSrv *http.Server, cfg *config.Config, _ user.Module, _ *slog.Logger, _ *goredis.Client, _ *amqp091.Connection, _ *pkgamqp.Publisher, grpcSrv *gogrpc.Server, api gohuma.API, consumerRunner consumer.Runner) *app.App {
	consumers := []func(ctx context.Context) error{(func(ctx context.Context) error)(consumerRunner)}
	return app.New(httpSrv, cfg, grpcSrv, api, consumers)
}

func InitializeApp(ctx context.Context) *app.App {
	wire.Build(
		config.SetupConfig,
		logger.SetupLogger,
		wire.FieldsOf(new(*config.Config), "App", "Logger", "DB", "JWT", "Redis", "GRPC", "AMQP"),
		db.Setup,
		redis.Setup,
		pkgamqp.Setup,
		pkgamqp.NewPublisher,
		event.NewEventBus,
		server.SetupMux,
		server.SetupHTTPServer,
		huma.Setup,
		pkggrpc.Setup,
		sharedjwt.NewJWTManager,
		persistence.NewUserRepository,
		service.NewUserLoaderCreator,
		middleware.Setup,

		consumer.SetupConsumers,

		user.InitializeUserModule,

		newApp,
	)
	return nil
}
