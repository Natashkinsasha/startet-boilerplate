package consumer

import (
	pkgamqp "starter-boilerplate/pkg/amqp"
)

type Init struct{}

func SetupConsumers(b *pkgamqp.Broker, profileCreated *ProfileCreatedConsumer) Init {
	profileCreated.Register(b)
	return Init{}
}
