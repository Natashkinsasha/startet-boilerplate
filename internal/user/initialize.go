//go:build wireinject

package user

import (
	sharedmw "starter-boilerplate/internal/shared/middleware"
	"starter-boilerplate/internal/user/app/usecase"
	"starter-boilerplate/internal/user/domain/repository"
	"starter-boilerplate/internal/user/transport/consumer"
	usercontract "starter-boilerplate/internal/user/transport/contract"
	"starter-boilerplate/internal/user/transport/handler"
	pkgamqp "starter-boilerplate/pkg/amqp"
	pkgjwt "starter-boilerplate/pkg/jwt"
	"starter-boilerplate/pkg/outbox"

	"starter-boilerplate/internal/user/app/service"

	"github.com/danielgtaylor/huma/v2"
	"github.com/google/wire"
	"github.com/uptrace/bun"
	gogrpc "google.golang.org/grpc"
)

type Module struct{}

func NewModule(_ handler.HandlersInit, _ usercontract.Init, _ consumer.Init) Module {
	return Module{}
}

func InitializeUserModule(_ *bun.DB, api huma.API, grpcSrv *gogrpc.Server, _ *pkgjwt.Manager, _ sharedmw.Init, _ repository.UserRepository, _ repository.ProfileRepository, _ outbox.Bus, _ *pkgamqp.Broker) Module {
	wire.Build(
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
		NewModule,
	)
	return Module{}
}
