//go:build unit

package grpc

import (
	"context"
	"errors"
	"net/http"
	"testing"

	"starter-boilerplate/pkg/apperror"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestHttpToGRPC(t *testing.T) {
	tests := []struct {
		httpStatus int
		expected   codes.Code
	}{
		{http.StatusBadRequest, codes.InvalidArgument},
		{http.StatusUnauthorized, codes.Unauthenticated},
		{http.StatusForbidden, codes.PermissionDenied},
		{http.StatusNotFound, codes.NotFound},
		{http.StatusConflict, codes.AlreadyExists},
		{http.StatusUnprocessableEntity, codes.InvalidArgument},
		{http.StatusTooManyRequests, codes.ResourceExhausted},
		{http.StatusInternalServerError, codes.Internal},
		{http.StatusServiceUnavailable, codes.Internal},
	}

	for _, tt := range tests {
		t.Run(http.StatusText(tt.httpStatus), func(t *testing.T) {
			assert.Equal(t, tt.expected, httpToGRPC(tt.httpStatus))
		})
	}
}

func TestErrorInterceptor_NoError(t *testing.T) {
	interceptor := ErrorInterceptor()
	handler := func(_ context.Context, _ any) (any, error) {
		return "ok", nil
	}
	info := &grpc.UnaryServerInfo{FullMethod: "/test/Method"}

	resp, err := interceptor(context.Background(), nil, info, handler)

	require.NoError(t, err)
	assert.Equal(t, "ok", resp)
}

func TestErrorInterceptor_AppError(t *testing.T) {
	interceptor := ErrorInterceptor()
	handler := func(_ context.Context, _ any) (any, error) {
		return nil, apperror.New(http.StatusNotFound, "user not found")
	}
	info := &grpc.UnaryServerInfo{FullMethod: "/test/Method"}

	resp, err := interceptor(context.Background(), nil, info, handler)

	assert.Nil(t, resp)
	require.Error(t, err)

	st, ok := status.FromError(err)
	require.True(t, ok)
	assert.Equal(t, codes.NotFound, st.Code())
	assert.Equal(t, "user not found", st.Message())
}

func TestErrorInterceptor_GenericError(t *testing.T) {
	interceptor := ErrorInterceptor()
	handler := func(_ context.Context, _ any) (any, error) {
		return nil, errors.New("unexpected failure")
	}
	info := &grpc.UnaryServerInfo{FullMethod: "/test/Method"}

	resp, err := interceptor(context.Background(), nil, info, handler)

	assert.Nil(t, resp)
	require.Error(t, err)

	st, ok := status.FromError(err)
	require.True(t, ok)
	assert.Equal(t, codes.Internal, st.Code())
	assert.Equal(t, "internal server error", st.Message())
}
