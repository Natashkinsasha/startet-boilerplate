package consumer

import (
	"context"
	"log/slog"

	sharedevent "starter-boilerplate/internal/shared/event"
	"starter-boilerplate/internal/user/app/service"
	pkgamqp "starter-boilerplate/pkg/amqp"
	"starter-boilerplate/pkg/event"

	amqp091 "github.com/rabbitmq/amqp091-go"
)

const queueProfileUpdater = "tag.profile"

type ProfileUpdaterConsumer struct {
	profileSvc *service.ProfileService
}

func NewProfileUpdaterConsumer(ps *service.ProfileService) *ProfileUpdaterConsumer {
	return &ProfileUpdaterConsumer{profileSvc: ps}
}

func (c *ProfileUpdaterConsumer) Register(b *pkgamqp.Broker) {
	r := sharedevent.NewRouter()
	sharedevent.Route(r, c.profileSvc.OnUserCreated)
	sharedevent.Route(r, c.profileSvc.OnPasswordChanged)
	r.Default(func(_ context.Context, _ []byte, meta pkgamqp.DeliveryMeta) error {
		slog.Warn("profile updater: unhandled event", slog.String("routing_key", meta.RoutingKey))
		return nil
	})
	pkgamqp.AddRawConsumer(b, pkgamqp.ConsumerConfig{
		Queue:    queueProfileUpdater,
		Exchange: event.ExchangeTagged,
		BindingArgs: amqp091.Table{
			"x-match":     "any",
			"tag.profile": true,
		},
	}, r.Handler())
}
