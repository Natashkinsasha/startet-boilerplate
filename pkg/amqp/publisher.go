package amqp

import (
	"context"
	"encoding/json"
	"fmt"

	amqp091 "github.com/rabbitmq/amqp091-go"
)

type Publisher struct {
	conn *amqp091.Connection
}

func NewPublisher(conn *amqp091.Connection) *Publisher {
	return &Publisher{conn: conn}
}

func (p *Publisher) Publish(ctx context.Context, exchange, routingKey string, body []byte) error {
	if p.conn == nil {
		return fmt.Errorf("amqp: connection is nil (standalone mode?)")
	}

	ch, err := p.conn.Channel()
	if err != nil {
		return fmt.Errorf("amqp: open channel: %w", err)
	}
	defer ch.Close()

	return ch.PublishWithContext(ctx, exchange, routingKey, false, false, amqp091.Publishing{
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
