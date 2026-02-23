package amqp

import (
	"context"
	"log/slog"
	"time"

	amqp091 "github.com/rabbitmq/amqp091-go"
)

// WithLogging returns a middleware that logs routing key, duration, and errors.
func WithLogging() Middleware {
	return func(next HandlerFunc) HandlerFunc {
		return func(ctx context.Context, msg amqp091.Delivery) error {
			start := time.Now()
			err := next(ctx, msg)
			attrs := []any{
				slog.String("routing_key", msg.RoutingKey),
				slog.Duration("duration", time.Since(start)),
			}
			if err != nil {
				slog.Error("message handling failed", append(attrs, slog.Any("error", err))...)
			} else {
				slog.Info("message handled", attrs...)
			}
			return err
		}
	}
}
