package consumer

import (
	amqp091 "github.com/rabbitmq/amqp091-go"

	pkgamqp "starter-boilerplate/pkg/amqp"
)

func Setup(conn *amqp091.Connection) *pkgamqp.Broker {
	b := pkgamqp.NewBroker(conn)
	b.Use(pkgamqp.WithRecover(), pkgamqp.WithLogging())
	return b
}
