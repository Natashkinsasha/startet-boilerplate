package amqp

import (
	"log/slog"

	amqp091 "github.com/rabbitmq/amqp091-go"
)

// Setup creates a new AMQP connection.
// *slog.Logger parameter ensures Wire initializes the logger before AMQP.
func Setup(cfg AMQPConfig, _ *slog.Logger) *amqp091.Connection {
	if cfg.Standalone {
		slog.Warn("standalone mode: skipping amqp connection")
		return nil
	}

	conn, err := amqp091.Dial(cfg.URL)
	if err != nil {
		panic("failed to connect to amqp: " + err.Error())
	}

	slog.Info("amqp connected", slog.String("url", cfg.URL))

	return conn
}
