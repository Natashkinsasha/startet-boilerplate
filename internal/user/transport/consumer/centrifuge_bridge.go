package consumer

import (
	"context"
	"encoding/json"
	"log/slog"

	"starter-boilerplate/internal/shared/centrifugenode"
	sharedevent "starter-boilerplate/internal/shared/event"
	userevent "starter-boilerplate/internal/user/domain/event"
	pkgamqp "starter-boilerplate/pkg/amqp"
	"starter-boilerplate/pkg/event"
)

const queueCentrifugeBridge = "centrifuge.bridge"

type BridgeConsumer struct {
	publisherSvc *centrifugenode.Publisher
}

func NewBridgeConsumer(p *centrifugenode.Publisher) *BridgeConsumer {
	return &BridgeConsumer{publisherSvc: p}
}

type BridgeInit struct{}

// SetupBridgeConsumer registers the bridge consumer on the AMQP broker.
// No-op if the publisher's node is nil (standalone mode).
func SetupBridgeConsumer(broker *pkgamqp.Broker, c *BridgeConsumer) BridgeInit {
	if !c.publisherSvc.Enabled() {
		return BridgeInit{}
	}
	c.Register(broker)
	return BridgeInit{}
}

func (c *BridgeConsumer) Register(b *pkgamqp.Broker) {
	r := sharedevent.NewRouter()
	sharedevent.Route(r, c.onUserCreated)
	sharedevent.Route(r, c.onPasswordChanged)
	sharedevent.Route(r, c.onUserLoggedIn)
	r.Default(func(_ context.Context, _ []byte, meta pkgamqp.DeliveryMeta) error {
		slog.Warn("centrifuge bridge: unhandled event", slog.String("routing_key", meta.RoutingKey))
		return nil
	})

	pkgamqp.AddRawConsumer(b, pkgamqp.ConsumerConfig{
		Queue:      queueCentrifugeBridge,
		Exchange:   event.ExchangeEvents,
		RoutingKey: "user.#",
	}, r.Handler())
}

func (c *BridgeConsumer) onUserCreated(ctx context.Context, e userevent.UserCreatedEvent, _ pkgamqp.DeliveryMeta) error {
	payload, _ := json.Marshal(e)
	return c.publisherSvc.PublishPersonal(ctx, e.UserID, userevent.UserCreated, payload)
}

func (c *BridgeConsumer) onPasswordChanged(ctx context.Context, e userevent.PasswordChangedEvent, _ pkgamqp.DeliveryMeta) error {
	payload, _ := json.Marshal(e)
	return c.publisherSvc.PublishPersonal(ctx, e.UserID, userevent.PasswordChanged, payload)
}

func (c *BridgeConsumer) onUserLoggedIn(ctx context.Context, e userevent.UserLoggedInEvent, _ pkgamqp.DeliveryMeta) error {
	payload, _ := json.Marshal(e)
	return c.publisherSvc.PublishPersonal(ctx, e.UserID, userevent.UserLoggedIn, payload)
}
