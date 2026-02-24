//go:build unit

package apperror

import (
	"errors"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) {
	err := New(http.StatusNotFound, "user not found")

	assert.Equal(t, http.StatusNotFound, err.Status)
	assert.Equal(t, "user not found", err.Message)
}

func TestAppError_Error(t *testing.T) {
	err := New(http.StatusBadRequest, "bad request")

	assert.Equal(t, "bad request", err.Error())
}

func TestAppError_GetStatus(t *testing.T) {
	err := New(http.StatusConflict, "conflict")

	assert.Equal(t, http.StatusConflict, err.GetStatus())
}

func TestWrap(t *testing.T) {
	cause := errors.New("connection refused")
	err := Wrap(cause, http.StatusInternalServerError, "db error")

	assert.Equal(t, http.StatusInternalServerError, err.Status)
	assert.Equal(t, "db error: connection refused", err.Message)
	assert.Equal(t, "db error: connection refused", err.Error())
}

func TestAppError_ImplementsError(t *testing.T) {
	var err error = New(http.StatusBadRequest, "test")

	var appErr *AppError
	require.True(t, errors.As(err, &appErr))
	assert.Equal(t, http.StatusBadRequest, appErr.GetStatus())
}
