package consumer

import (
	pkgamqp "starter-boilerplate/pkg/amqp"
)

type Init struct{}

func SetupConsumers(b *pkgamqp.Broker, profileUpdater *ProfileUpdaterConsumer) Init {
	profileUpdater.Register(b)
	return Init{}
}
