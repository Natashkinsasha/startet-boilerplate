package event

import (
	"context"

	pkgamqp "starter-boilerplate/pkg/amqp"
	"starter-boilerplate/pkg/outbox"

	amqp091 "github.com/rabbitmq/amqp091-go"
)

// OutboxPublisher adapts the outbox.Publisher interface to AMQP via Broker.
type OutboxPublisher struct {
	broker   *pkgamqp.Broker
	exchange string
}

func NewOutboxPublisher(broker *pkgamqp.Broker, exchange string) *OutboxPublisher {
	return &OutboxPublisher{broker: broker, exchange: exchange}
}

// NewDefaultOutboxPublisher is a Wire-friendly constructor that hardcodes ExchangeEvents.
// The Bus parameter is unused but ensures Wire initializes the bus first.
func NewDefaultOutboxPublisher(_ Bus, broker *pkgamqp.Broker) *OutboxPublisher {
	return NewOutboxPublisher(broker, ExchangeEvents)
}

func (p *OutboxPublisher) Publish(ctx context.Context, entry outbox.Entry) error {
	return p.broker.Publish(ctx, p.exchange, entry.EventName, toAMQPTable(entry.Headers), entry.Payload, pkgamqp.AtLeastOnce)
}

func toAMQPTable(headers map[string]any) amqp091.Table {
	if len(headers) == 0 {
		return nil
	}
	t := make(amqp091.Table, len(headers))
	for k, v := range headers {
		t[k] = v
	}
	return t
}
