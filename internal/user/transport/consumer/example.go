package consumer

import (
	"context"
	"log/slog"

	amqp091 "github.com/rabbitmq/amqp091-go"

	pkgamqp "starter-boilerplate/pkg/amqp"
)

func newExampleConsumer() *pkgamqp.Consumer {
	return pkgamqp.NewConsumer(
		pkgamqp.ConsumerConfig{Queue: "user.example"},
		handleExampleMessage,
	)
}

func handleExampleMessage(_ context.Context, msg amqp091.Delivery) error {
	slog.Info("received message", slog.String("body", string(msg.Body)))
	return nil
}
