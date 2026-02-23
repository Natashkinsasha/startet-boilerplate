package amqp

import (
	"context"
	"encoding/json"
	"fmt"

	amqp091 "github.com/rabbitmq/amqp091-go"
)

// publisher publishes messages using a single long-lived AMQP channel.
// Internal to [Broker] — use Broker.Publish / Broker.PublishJSON.
type publisher struct {
	ch *amqp091.Channel
}

func newPublisher(conn *amqp091.Connection) *publisher {
	if conn == nil {
		return &publisher{}
	}

	ch, err := conn.Channel()
	if err != nil {
		panic(fmt.Sprintf("amqp: open publisher channel: %v", err))
	}

	return &publisher{ch: ch}
}

func (p *publisher) Close() error {
	if p.ch == nil {
		return nil
	}
	return p.ch.Close()
}

func (p *publisher) Publish(ctx context.Context, exchange, routingKey string, body []byte) error {
	if p.ch == nil {
		return fmt.Errorf("amqp: channel is nil (standalone mode?)")
	}

	return p.ch.PublishWithContext(ctx, exchange, routingKey, false, false, amqp091.Publishing{
		ContentType: "application/json",
		Body:        body,
	})
}

func (p *publisher) PublishJSON(ctx context.Context, exchange, routingKey string, payload any) error {
	if err := validate.StructCtx(ctx, payload); err != nil {
		return fmt.Errorf("amqp: validate: %w", err)
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("amqp: marshal: %w", err)
	}

	return p.Publish(ctx, exchange, routingKey, body)
}
