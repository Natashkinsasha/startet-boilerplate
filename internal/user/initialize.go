//go:build wireinject

package user

import (
	sharedmw "starter-boilerplate/internal/shared/middleware"
	"starter-boilerplate/internal/user/app/service"
	"starter-boilerplate/internal/user/app/usecase"
	"starter-boilerplate/internal/user/infra/persistence"
	usercontract "starter-boilerplate/internal/user/transport/contract"
	"starter-boilerplate/internal/user/transport/handler"
	pkgjwt "starter-boilerplate/pkg/jwt"

	"github.com/danielgtaylor/huma/v2"
	"github.com/google/wire"
	"github.com/uptrace/bun"
	gogrpc "google.golang.org/grpc"
)

type Module struct{}

func NewModule(_ handler.HandlersInit, _ usercontract.Init) Module {
	return Module{}
}

func InitializeUserModule(db *bun.DB, api huma.API, grpcSrv *gogrpc.Server, jwtManager *pkgjwt.Manager, _ sharedmw.Init) Module {
	wire.Build(
		persistence.NewUserRepository,
		service.NewUserService,
		service.NewTokenService,
		usecase.NewLoginUseCase,
		usecase.NewRefreshUseCase,
		handler.NewLoginHandler,
		handler.NewRefreshHandler,
		handler.NewGetUserHandler,
		handler.SetupHandlers,
		usercontract.SetupUserContract,
		NewModule,
	)
	return Module{}
}
