package event

import (
	"context"
	"fmt"

	pkgamqp "starter-boilerplate/pkg/amqp"

	amqp091 "github.com/rabbitmq/amqp091-go"
)

const (
	ExchangeEvents = "events"
	ExchangeDLX    = "dlx"
)

// AMQPBus implements Bus by publishing events to a single topic exchange.
// The routing key is derived from Event.EventName().
type AMQPBus struct {
	broker   *pkgamqp.Broker
	exchange string
}

// AMQPBusOption configures an AMQPBus.
type AMQPBusOption func(*AMQPBus)

func NewAMQPBus(broker *pkgamqp.Broker, exchange string, opts ...AMQPBusOption) *AMQPBus {
	b := &AMQPBus{broker: broker, exchange: exchange}
	for _, o := range opts {
		o(b)
	}
	return b
}

// NewEventBus is a Wire provider that creates a Bus backed by AMQP
// with the "events" topic exchange.
// It declares the exchange on startup so both publishers and consumers
// can rely on it regardless of initialization order.
func NewEventBus(conn *amqp091.Connection, broker *pkgamqp.Broker) Bus {
	if conn != nil {
		if err := declareExchange(conn, ExchangeEvents); err != nil {
			panic(fmt.Sprintf("event bus: declare exchange: %v", err))
		}
	}
	return NewAMQPBus(broker, ExchangeEvents)
}

func declareExchange(conn *amqp091.Connection, name string) error {
	ch, err := conn.Channel()
	if err != nil {
		return err
	}
	defer func() { _ = ch.Close() }()
	return ch.ExchangeDeclare(name, "topic", true, false, false, false, nil)
}

func (b *AMQPBus) Publish(ctx context.Context, e Event) error {
	return b.broker.PublishJSON(ctx, b.exchange, e.EventName(), e, pkgamqp.AtLeastOnce)
}
