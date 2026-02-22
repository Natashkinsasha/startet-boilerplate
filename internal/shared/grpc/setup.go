package grpc

import (
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

type GRPCConfig struct {
	Port int `yaml:"port" validate:"required"`
}

// Setup creates a new gRPC server.
// *zap.Logger parameter ensures Wire initializes the logger before gRPC.
func Setup(cfg GRPCConfig, _ *zap.Logger) *grpc.Server {
	zap.L().Info("grpc server created", zap.Int("port", cfg.Port))
	return grpc.NewServer(
		grpc.UnaryInterceptor(ErrorInterceptor()),
	)
}
