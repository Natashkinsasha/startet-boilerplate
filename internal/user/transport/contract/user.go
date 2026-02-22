package contract

import (
	"context"

	gen "starter-boilerplate/gen/user"
	"starter-boilerplate/internal/user/domain/repository"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Contract struct {
	gen.UnimplementedUserContractServer
	ur repository.UserRepository
}

type Init struct{}

func SetupUserContract(grpcSrv *grpc.Server, ur repository.UserRepository) Init {
	c := &Contract{ur: ur}
	gen.RegisterUserContractServer(grpcSrv, c)
	return Init{}
}

func (c *Contract) GetUser(ctx context.Context, req *gen.GetUserRequest) (*gen.GetUserResponse, error) {
	u, err := c.ur.FindByID(ctx, req.Id)
	if err != nil {
		return nil, status.Error(codes.NotFound, "user not found")
	}
	if u == nil {
		return nil, status.Error(codes.NotFound, "user not found")
	}
	return &gen.GetUserResponse{Id: u.ID, Email: u.Email}, nil
}
