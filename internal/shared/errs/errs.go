package errs

import (
	"net/http"

	"starter-boilerplate/pkg/apperror"
)

var (
	ErrAccessDenied       = apperror.New(http.StatusForbidden, "access denied")
	ErrNotFound           = apperror.New(http.StatusNotFound, "not found")
	ErrInvalidCredentials = apperror.New(http.StatusUnauthorized, "invalid credentials")
	ErrInvalidToken       = apperror.New(http.StatusUnauthorized, "invalid token")
	ErrEmailAlreadyExists = apperror.New(http.StatusConflict, "email already exists")
)
