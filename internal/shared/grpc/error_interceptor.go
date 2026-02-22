package grpc

import (
	"context"
	"errors"
	"net/http"

	apperror "starter-boilerplate/internal/shared/error"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func ErrorInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
		resp, err := handler(ctx, req)
		if err == nil {
			return resp, nil
		}

		var appErr *apperror.AppError
		if errors.As(err, &appErr) {
			return nil, status.Error(httpToGRPC(appErr.Status), appErr.Message)
		}

		zap.L().Error("internal error",
			zap.String("method", info.FullMethod),
			zap.Error(err),
		)
		return nil, status.Error(codes.Internal, "internal server error")
	}
}

func httpToGRPC(httpStatus int) codes.Code {
	switch httpStatus {
	case http.StatusBadRequest:
		return codes.InvalidArgument
	case http.StatusUnauthorized:
		return codes.Unauthenticated
	case http.StatusForbidden:
		return codes.PermissionDenied
	case http.StatusNotFound:
		return codes.NotFound
	case http.StatusConflict:
		return codes.AlreadyExists
	case http.StatusUnprocessableEntity:
		return codes.InvalidArgument
	case http.StatusTooManyRequests:
		return codes.ResourceExhausted
	default:
		return codes.Internal
	}
}
