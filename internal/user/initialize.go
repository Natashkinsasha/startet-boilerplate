//go:build wireinject

package user

import (
	"starter-boilerplate/internal/shared/centrifugenode"
	"starter-boilerplate/internal/shared/middleware"
	"starter-boilerplate/internal/user/app/service"
	"starter-boilerplate/internal/user/app/usecase"
	"starter-boilerplate/internal/user/infra/persistence"
	"starter-boilerplate/internal/user/transport/consumer"
	usercontract "starter-boilerplate/internal/user/transport/contract"
	"starter-boilerplate/internal/user/transport/handler"
	pkgamqp "starter-boilerplate/pkg/amqp"
	pkgdb "starter-boilerplate/pkg/db"
	pkgjwt "starter-boilerplate/pkg/jwt"
	"starter-boilerplate/pkg/outbox"

	"github.com/danielgtaylor/huma/v2"
	"github.com/google/wire"
	"github.com/uptrace/bun"
	gogrpc "google.golang.org/grpc"
)

type Module struct{}

func NewModule(_ handler.HandlersInit, _ usercontract.Init, _ consumer.Init, _ consumer.BridgeInit) Module {
	return Module{}
}

func InitializeUserModule(api huma.API, grpcSrv *gogrpc.Server, _ *pkgjwt.Manager, _ *bun.DB, _ outbox.Bus, _ *pkgamqp.Broker, _ pkgdb.UoW, _ *centrifugenode.Publisher, _ middleware.Init) Module {
	wire.Build(
		persistence.NewUserRepository,
		persistence.NewProfileRepository,
		service.NewUserService,
		service.NewTokenService,
		usecase.NewLoginUseCase,
		usecase.NewRefreshUseCase,
		usecase.NewGetUserUseCase,
		usecase.NewRegisterUseCase,
		usecase.NewChangePasswordUseCase,
		service.NewProfileService,
		handler.NewLoginHandler,
		handler.NewRefreshHandler,
		handler.NewGetUserHandler,
		handler.NewRegisterHandler,
		handler.NewChangePasswordHandler,
		handler.SetupHandlers,
		usercontract.SetupUserContract,
		consumer.NewProfileUpdaterConsumer,
		consumer.SetupConsumers,
		consumer.NewBridgeConsumer,
		consumer.SetupBridgeConsumer,
		NewModule,
	)
	return Module{}
}
