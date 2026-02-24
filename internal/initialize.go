//go:build wireinject

package internal

import (
	"context"
	"log/slog"
	"net/http"

	"starter-boilerplate/internal/shared/app"
	"starter-boilerplate/internal/shared/config"
	sharedconsumer "starter-boilerplate/internal/shared/consumer"
	"starter-boilerplate/internal/shared/huma"
	sharedjwt "starter-boilerplate/internal/shared/jwt"
	"starter-boilerplate/internal/shared/logger"
	"starter-boilerplate/internal/shared/middleware"
	"starter-boilerplate/internal/shared/server"
	"starter-boilerplate/internal/user"
	"starter-boilerplate/internal/user/app/service"
	"starter-boilerplate/internal/user/infra/persistence"
	pkgamqp "starter-boilerplate/pkg/amqp"
	"starter-boilerplate/pkg/db"
	"starter-boilerplate/pkg/event"
	pkggrpc "starter-boilerplate/pkg/grpc"
	"starter-boilerplate/pkg/outbox"
	"starter-boilerplate/pkg/redis"

	gohuma "github.com/danielgtaylor/huma/v2"
	"github.com/google/wire"
	goredis "github.com/redis/go-redis/v9"
	gogrpc "google.golang.org/grpc"
)

func newApp(httpSrv *http.Server, cfg *config.Config, _ user.Module, _ *slog.Logger, _ *goredis.Client, grpcSrv *gogrpc.Server, api gohuma.API, broker *pkgamqp.Broker, relay *outbox.Relay) *app.App {
	return app.New(httpSrv, cfg, grpcSrv, api, broker, relay)
}

func provideOutboxPublisher(broker *pkgamqp.Broker, _ event.Bus) *event.OutboxPublisher {
	return event.NewOutboxPublisher(broker, event.ExchangeEvents)
}

func provideRelayConfig(cfg *config.Config) outbox.RelayConfig {
	return cfg.Outbox
}

func InitializeApp(ctx context.Context) *app.App {
	wire.Build(
		config.SetupConfig,
		logger.SetupLogger,
		wire.FieldsOf(new(*config.Config), "App", "Logger", "DB", "JWT", "Redis", "GRPC", "AMQP"),
		db.Setup,
		redis.Setup,
		pkgamqp.Setup,
		event.NewEventBus,
		server.SetupMux,
		server.SetupHTTPServer,
		huma.Setup,
		pkggrpc.Setup,
		sharedjwt.NewJWTManager,
		persistence.NewUserRepository,
		persistence.NewProfileRepository,
		service.NewUserLoaderCreator,
		middleware.Setup,

		outbox.NewRepository,
		outbox.NewOutboxBus,
		wire.Bind(new(outbox.Bus), new(*outbox.OutboxBus)),
		provideOutboxPublisher,
		wire.Bind(new(outbox.Publisher), new(*event.OutboxPublisher)),
		provideRelayConfig,
		outbox.NewRelay,

		sharedconsumer.Setup,
		user.InitializeUserModule,

		newApp,
	)
	return nil
}
