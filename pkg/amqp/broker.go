package amqp

import (
	"context"
	"log/slog"

	amqp091 "github.com/rabbitmq/amqp091-go"
)

// Broker is the top-level AMQP component that owns the Publisher
// and ConsumerGroup. It provides a single Run / Shutdown lifecycle.
type Broker struct {
	Publisher *Publisher
	Consumers *ConsumerGroup
}

// NewBroker creates a Broker with a Publisher and ConsumerGroup
// backed by the given connection. If conn is nil (standalone mode),
// both sub-components handle it gracefully.
func NewBroker(conn *amqp091.Connection) *Broker {
	return &Broker{
		Publisher: NewPublisher(conn),
		Consumers: NewConsumerGroup(conn),
	}
}

// Use sets group-level middlewares for all consumers.
func (b *Broker) Use(mws ...Middleware) {
	b.Consumers.Use(mws...)
}

// Run starts all consumers and blocks until ctx is done.
func (b *Broker) Run(ctx context.Context) error {
	return b.Consumers.Run(ctx)
}

// Shutdown gracefully stops consumers and closes the publisher channel.
func (b *Broker) Shutdown() {
	b.Consumers.Shutdown()

	if err := b.Publisher.Close(); err != nil {
		slog.Error("failed to close publisher", slog.Any("error", err))
	}
}
