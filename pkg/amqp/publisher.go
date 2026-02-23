package amqp

import (
	"context"
	"encoding/json"
	"fmt"

	amqp091 "github.com/rabbitmq/amqp091-go"
)

// Publisher publishes messages to AMQP exchanges using a single long-lived channel.
type Publisher struct {
	ch *amqp091.Channel
}

// NewPublisher opens a channel on the given connection.
// If conn is nil (standalone mode), Publish will return an error.
func NewPublisher(conn *amqp091.Connection) *Publisher {
	if conn == nil {
		return &Publisher{}
	}

	ch, err := conn.Channel()
	if err != nil {
		panic(fmt.Sprintf("amqp: open publisher channel: %v", err))
	}

	return &Publisher{ch: ch}
}

// Close closes the underlying AMQP channel.
func (p *Publisher) Close() error {
	if p.ch == nil {
		return nil
	}
	return p.ch.Close()
}

func (p *Publisher) Publish(ctx context.Context, exchange, routingKey string, body []byte) error {
	if p.ch == nil {
		return fmt.Errorf("amqp: channel is nil (standalone mode?)")
	}

	return p.ch.PublishWithContext(ctx, exchange, routingKey, false, false, amqp091.Publishing{
		ContentType: "application/json",
		Body:        body,
	})
}

// PublishJSON validates, marshals the payload to JSON, and publishes it.
func (p *Publisher) PublishJSON(ctx context.Context, exchange, routingKey string, payload any) error {
	if err := validate.StructCtx(ctx, payload); err != nil {
		return fmt.Errorf("amqp: validate: %w", err)
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("amqp: marshal: %w", err)
	}

	return p.Publish(ctx, exchange, routingKey, body)
}
