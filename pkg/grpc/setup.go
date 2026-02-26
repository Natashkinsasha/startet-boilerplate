package grpc

import (
	"log/slog"

	"google.golang.org/grpc"
)

type GRPCConfig struct {
	Port int `yaml:"port" validate:"required"`
}

// Setup creates a new gRPC server with the default error interceptor.
// *slog.Logger parameter ensures Wire initializes the logger before gRPC.
func Setup(cfg GRPCConfig, _ *slog.Logger) *grpc.Server {
	slog.Info("grpc server created", slog.Int("port", cfg.Port))
	return grpc.NewServer(
		grpc.ChainUnaryInterceptor(ErrorInterceptor()),
	)
}
