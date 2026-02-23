package event

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"

	pkgamqp "starter-boilerplate/pkg/amqp"
	pkgevent "starter-boilerplate/pkg/event"
)

// Router dispatches raw messages to typed handlers by routing key.
type Router struct {
	routes   map[string]func(ctx context.Context, body []byte, meta pkgamqp.DeliveryMeta) error
	fallback func(ctx context.Context, body []byte, meta pkgamqp.DeliveryMeta) error
}

// NewRouter creates an empty router.
func NewRouter() *Router {
	return &Router{
		routes: make(map[string]func(ctx context.Context, body []byte, meta pkgamqp.DeliveryMeta) error),
	}
}

// Route registers a typed handler keyed by T's EventName().
// Unmarshal and validation happen automatically before fn is called.
func Route[T pkgevent.Event](r *Router, fn func(ctx context.Context, payload T, meta pkgamqp.DeliveryMeta) error) {
	var zero T
	r.routes[zero.EventName()] = func(ctx context.Context, body []byte, meta pkgamqp.DeliveryMeta) error {
		var payload T
		if err := json.Unmarshal(body, &payload); err != nil {
			return fmt.Errorf("unmarshal: %w", err)
		}
		if err := pkgamqp.Validate(ctx, payload); err != nil {
			return fmt.Errorf("validate: %w", err)
		}
		return fn(ctx, payload, meta)
	}
}

// Default sets a fallback handler for unmatched routing keys.
func (r *Router) Default(fn func(ctx context.Context, body []byte, meta pkgamqp.DeliveryMeta) error) {
	r.fallback = fn
}

// Handler returns a function compatible with pkgamqp.AddRawConsumer.
func (r *Router) Handler() func(ctx context.Context, body []byte, meta pkgamqp.DeliveryMeta) error {
	return func(ctx context.Context, body []byte, meta pkgamqp.DeliveryMeta) error {
		if h, ok := r.routes[meta.RoutingKey]; ok {
			return h(ctx, body, meta)
		}
		if r.fallback != nil {
			return r.fallback(ctx, body, meta)
		}
		slog.Warn("event router: unhandled routing key", slog.String("routing_key", meta.RoutingKey))
		return nil
	}
}
