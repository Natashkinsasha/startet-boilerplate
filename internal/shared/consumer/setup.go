package consumer

import (
	pkgamqp "starter-boilerplate/pkg/amqp"

	amqp091 "github.com/rabbitmq/amqp091-go"
)

func Setup(conn *amqp091.Connection, cfg pkgamqp.AMQPConfig) *pkgamqp.Broker {
	b := pkgamqp.NewBroker(conn, cfg.Pool)
	b.Use(pkgamqp.WithRecover(), pkgamqp.WithLogging())
	return b
}
