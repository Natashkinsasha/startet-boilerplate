package consumer

import (
	amqp091 "github.com/rabbitmq/amqp091-go"

	pkgamqp "starter-boilerplate/pkg/amqp"
)

func Setup(conn *amqp091.Connection, cfg pkgamqp.AMQPConfig) *pkgamqp.Broker {
	b := pkgamqp.NewBroker(conn, cfg.Pool)
	b.Use(pkgamqp.WithRecover(), pkgamqp.WithLogging())
	return b
}
