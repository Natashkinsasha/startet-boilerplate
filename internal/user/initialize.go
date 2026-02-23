//go:build wireinject

package user

import (
	sharedmw "starter-boilerplate/internal/shared/middleware"
	"starter-boilerplate/internal/user/app/usecase"
	"starter-boilerplate/internal/user/domain/repository"
	usercontract "starter-boilerplate/internal/user/transport/contract"
	"starter-boilerplate/internal/user/transport/handler"
	"starter-boilerplate/pkg/event"
	pkgjwt "starter-boilerplate/pkg/jwt"

	"starter-boilerplate/internal/user/app/service"

	"github.com/danielgtaylor/huma/v2"
	"github.com/google/wire"
	"github.com/uptrace/bun"
	gogrpc "google.golang.org/grpc"
)

type Module struct{}

func NewModule(_ handler.HandlersInit, _ usercontract.Init) Module {
	return Module{}
}

func InitializeUserModule(_ *bun.DB, api huma.API, grpcSrv *gogrpc.Server, _ *pkgjwt.Manager, _ sharedmw.Init, _ repository.UserRepository, _ event.Bus) Module {
	wire.Build(
		service.NewUserService,
		service.NewTokenService,
		usecase.NewLoginUseCase,
		usecase.NewRefreshUseCase,
		usecase.NewGetUserUseCase,
		usecase.NewRegisterUseCase,
		handler.NewLoginHandler,
		handler.NewRefreshHandler,
		handler.NewGetUserHandler,
		handler.NewRegisterHandler,
		handler.SetupHandlers,
		usercontract.SetupUserContract,
		NewModule,
	)
	return Module{}
}
