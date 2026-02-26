//go:build wireinject

package internal

import (
	"context"
	"log/slog"
	"net/http"

	"starter-boilerplate/internal/shared/app"
	"starter-boilerplate/internal/shared/centrifugenode"
	"starter-boilerplate/internal/shared/config"
	sharedconsumer "starter-boilerplate/internal/shared/consumer"
	"starter-boilerplate/internal/shared/huma"
	sharedjwt "starter-boilerplate/internal/shared/jwt"
	"starter-boilerplate/internal/shared/logger"
	"starter-boilerplate/internal/shared/middleware"
	"starter-boilerplate/internal/shared/server"
	"starter-boilerplate/internal/user"
	pkgamqp "starter-boilerplate/pkg/amqp"
	pkgcentrifuge "starter-boilerplate/pkg/centrifuge"
	pkgdb "starter-boilerplate/pkg/db"
	"starter-boilerplate/pkg/event"
	pkggrpc "starter-boilerplate/pkg/grpc"
	"starter-boilerplate/pkg/outbox"
	"starter-boilerplate/pkg/redis"

	gocentrifuge "github.com/centrifugal/centrifuge"
	gohuma "github.com/danielgtaylor/huma/v2"
	"github.com/google/wire"
	goredis "github.com/redis/go-redis/v9"
	gogrpc "google.golang.org/grpc"
)

func newApp(httpSrv *http.Server, cfg *config.Config, _ user.Module, _ middleware.Init, _ *slog.Logger, _ *goredis.Client, grpcSrv *gogrpc.Server, api gohuma.API, broker *pkgamqp.Broker, relay *outbox.Relay, centrifugeNode *gocentrifuge.Node, _ centrifugenode.Init) *app.App {
	return app.New(httpSrv, cfg, grpcSrv, api, broker, relay, centrifugeNode)
}

func InitializeApp(ctx context.Context) *app.App {
	wire.Build(
		config.SetupConfig,
		logger.SetupLogger,
		wire.FieldsOf(new(*config.Config), "App", "Logger", "DB", "JWT", "Redis", "GRPC", "AMQP", "Outbox", "Centrifuge"),

		wire.NewSet(pkgdb.Setup, pkgdb.NewUnitOfWork, wire.Bind(new(pkgdb.UoW), new(*pkgdb.UnitOfWork))),
		redis.Setup,
		pkgamqp.Setup,
		wire.NewSet(server.SetupMux, server.SetupHTTPServer),
		huma.Setup,
		pkggrpc.Setup,
		sharedjwt.NewJWTManager,

		wire.NewSet(event.NewEventBus, event.NewDefaultOutboxPublisher, wire.Bind(new(outbox.Publisher), new(*event.OutboxPublisher))),
		wire.NewSet(outbox.NewRepository, outbox.NewOutboxBus, wire.Bind(new(outbox.Bus), new(*outbox.OutboxBus)), outbox.NewRelay),

		pkgcentrifuge.Setup,
		wire.NewSet(centrifugenode.NewPublisher, centrifugenode.Setup),

		middleware.Setup,
		sharedconsumer.Setup,
		user.InitializeUserModule,

		newApp,
	)
	return nil
}
