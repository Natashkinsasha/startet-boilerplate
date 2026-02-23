package event

import (
	"context"

	pkgamqp "starter-boilerplate/pkg/amqp"
)

// AMQPBus implements Bus by publishing events to a single topic exchange.
// The routing key is derived from Event.EventName().
type AMQPBus struct {
	publisher *pkgamqp.Publisher
	exchange  string
}

func NewAMQPBus(publisher *pkgamqp.Publisher, exchange string) *AMQPBus {
	return &AMQPBus{publisher: publisher, exchange: exchange}
}

// NewEventBus is a Wire provider that creates a Bus backed by AMQP
// with a default "events" topic exchange.
func NewEventBus(publisher *pkgamqp.Publisher) Bus {
	return NewAMQPBus(publisher, "events")
}

func (b *AMQPBus) Publish(ctx context.Context, e Event) error {
	return b.publisher.PublishJSON(ctx, b.exchange, e.EventName(), e)
}
