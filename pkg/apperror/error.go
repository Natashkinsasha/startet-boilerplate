package apperror

import "fmt"

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
