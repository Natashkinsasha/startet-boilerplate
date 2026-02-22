package apperror

import (
	"fmt"
	"net/http"
)

type AppError struct {
	Status  int
	Message string
}

func (e *AppError) Error() string {
	return e.Message
}

func (e *AppError) GetStatus() int {
	return e.Status
}

func New(status int, message string) *AppError {
	return &AppError{Status: status, Message: message}
}

func Wrap(err error, status int, message string) *AppError {
	return &AppError{Status: status, Message: fmt.Sprintf("%s: %v", message, err)}
}

var (
	ErrAccessDenied       = New(http.StatusForbidden, "access denied")
	ErrNotFound           = New(http.StatusNotFound, "not found")
	ErrInvalidCredentials = New(http.StatusUnauthorized, "invalid credentials")
	ErrInvalidToken       = New(http.StatusUnauthorized, "invalid token")
	ErrEmailAlreadyExists = New(http.StatusConflict, "email already exists")
)
