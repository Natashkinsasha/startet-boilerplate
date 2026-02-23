package amqp

import (
	"context"
	"fmt"

	amqp091 "github.com/rabbitmq/amqp091-go"
)

// publisher is the internal interface for message publishing.
type publisher interface {
	Publish(ctx context.Context, exchange, routingKey string, headers amqp091.Table, body []byte) error
	Close() error
}

// --- fire-and-forget (AtMostOnce) ----------------------------------------

// fireAndForgetPublisher publishes messages using a single long-lived AMQP
// channel with no delivery confirmation.
type fireAndForgetPublisher struct {
	ch *amqp091.Channel
}

func newSingleFireAndForget(conn *amqp091.Connection) *fireAndForgetPublisher {
	ch, err := conn.Channel()
	if err != nil {
		panic(fmt.Sprintf("amqp: open publisher channel: %v", err))
	}
	return &fireAndForgetPublisher{ch: ch}
}

func (p *fireAndForgetPublisher) Close() error {
	if p.ch == nil {
		return nil
	}
	return p.ch.Close()
}

func (p *fireAndForgetPublisher) Publish(ctx context.Context, exchange, routingKey string, headers amqp091.Table, body []byte) error {
	return p.ch.PublishWithContext(ctx, exchange, routingKey, false, false, amqp091.Publishing{
		Headers:     headers,
		ContentType: "application/json",
		Body:        body,
	})
}

// --- confirmed (AtLeastOnce) ---------------------------------------------

// confirmedPublisher opens a channel in confirm mode.
// Every publish waits for an ack from the broker and sets DeliveryMode to Persistent.
type confirmedPublisher struct {
	ch *amqp091.Channel
}

func newSingleConfirmed(conn *amqp091.Connection) *confirmedPublisher {
	ch, err := conn.Channel()
	if err != nil {
		panic(fmt.Sprintf("amqp: open confirmed publisher channel: %v", err))
	}

	if err := ch.Confirm(false); err != nil {
		_ = ch.Close()
		panic(fmt.Sprintf("amqp: enable confirm mode: %v", err))
	}

	return &confirmedPublisher{ch: ch}
}

func (p *confirmedPublisher) Close() error {
	if p.ch == nil {
		return nil
	}
	return p.ch.Close()
}

func (p *confirmedPublisher) Publish(ctx context.Context, exchange, routingKey string, headers amqp091.Table, body []byte) error {
	conf, err := p.ch.PublishWithDeferredConfirmWithContext(ctx, exchange, routingKey, false, false, amqp091.Publishing{
		Headers:      headers,
		ContentType:  "application/json",
		DeliveryMode: amqp091.Persistent,
		Body:         body,
	})
	if err != nil {
		return fmt.Errorf("amqp: confirmed publish: %w", err)
	}

	if !conf.Wait() {
		return fmt.Errorf("amqp: publish nacked by broker")
	}
	return nil
}

// --- factory -------------------------------------------------------------

func newSinglePublisher(conn *amqp091.Connection, g DeliveryGuarantee) publisher {
	switch g {
	case AtLeastOnce:
		return newSingleConfirmed(conn)
	default:
		return newSingleFireAndForget(conn)
	}
}
