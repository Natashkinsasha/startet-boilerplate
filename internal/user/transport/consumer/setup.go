package consumer

import (
	pkgamqp "starter-boilerplate/pkg/amqp"
)

type Init struct{}

func SetupConsumers(g *pkgamqp.ConsumerGroup, profileCreated *ProfileCreatedConsumer) Init {
	profileCreated.Register(g)
	return Init{}
}
