package amqp

import (
	"context"
	"fmt"
	"log/slog"
	"runtime/debug"

	amqp091 "github.com/rabbitmq/amqp091-go"
)

// WithRecover returns a middleware that catches panics and converts them to errors.
func WithRecover() Middleware {
	return func(next HandlerFunc) HandlerFunc {
		return func(ctx context.Context, msg amqp091.Delivery) (err error) {
			defer func() {
				if r := recover(); r != nil {
					slog.Error("consumer panic",
						slog.Any("panic", r),
						slog.String("stack", string(debug.Stack())),
					)
					err = fmt.Errorf("panic: %v", r)
				}
			}()
			return next(ctx, msg)
		}
	}
}
