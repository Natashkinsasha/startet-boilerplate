package amqp

import (
	"context"
	"log/slog"

	amqp091 "github.com/rabbitmq/amqp091-go"
)

// Broker is the single entry point for AMQP messaging.
// It encapsulates publishing and consuming behind one interface, manages
// the full lifecycle (Run / Shutdown), and owns the underlying connection resources.
//
// Publishing:
//
//	broker.Publish(ctx, exchange, key, body)
//	broker.PublishJSON(ctx, exchange, key, payload)  // with validation + marshal
//
// Consuming — register handlers before calling Run:
//
//	amqp.AddConsumer[EventT](broker, cfg, handler)
//	broker.Use(amqp.WithRecover(), amqp.WithLogging())  // group middlewares
//	broker.Run(ctx)  // blocks until ctx is done or a consumer fails
//
// Standalone mode: if conn is nil, Publish returns an error
// and Run blocks until ctx is done without starting consumers.
type Broker struct {
	publisher *publisher
	consumers *consumerGroup
}

// NewBroker creates a ready-to-use Broker backed by the given connection.
// Pass nil for standalone mode (no AMQP server required).
func NewBroker(conn *amqp091.Connection) *Broker {
	return &Broker{
		publisher: newPublisher(conn),
		consumers: newConsumerGroup(conn),
	}
}

// Publish sends a raw message to the given exchange with the specified routing key.
func (b *Broker) Publish(ctx context.Context, exchange, routingKey string, body []byte) error {
	return b.publisher.Publish(ctx, exchange, routingKey, body)
}

// PublishJSON validates the payload struct, marshals it to JSON, and publishes.
func (b *Broker) PublishJSON(ctx context.Context, exchange, routingKey string, payload any) error {
	return b.publisher.PublishJSON(ctx, exchange, routingKey, payload)
}

// Use appends group-level middlewares that apply to every registered consumer.
// Must be called before AddConsumer so that middlewares are captured at registration time.
func (b *Broker) Use(mws ...Middleware) {
	b.consumers.Use(mws...)
}

// Run declares all queues/bindings, then starts consuming on dedicated channels.
// It blocks until ctx is cancelled, Shutdown is called, or a consumer returns an error.
// On shutdown the delivery loop stops immediately, but in-flight handlers finish.
func (b *Broker) Run(ctx context.Context) error {
	return b.consumers.Run(ctx)
}

// Shutdown gracefully stops consumers (no new deliveries, in-flight handlers finish)
// and closes the publisher channel.
func (b *Broker) Shutdown() {
	b.consumers.Shutdown()

	if err := b.publisher.Close(); err != nil {
		slog.Error("failed to close publisher", slog.Any("error", err))
	}
}
