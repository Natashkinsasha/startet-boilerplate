package amqp

import (
	"context"
	"encoding/json"
	"fmt"

	amqp091 "github.com/rabbitmq/amqp091-go"
)

// Broker is the single entry point for AMQP messaging.
// It encapsulates publishing and consuming behind one interface, manages
// the full lifecycle (Run / Shutdown), and owns the underlying connection resources.
//
// Publishing:
//
//	broker.Publish(ctx, exchange, key, body)
//	broker.PublishJSON(ctx, exchange, key, payload)  // with validation + marshal
//
// Consuming — register handlers before calling Run:
//
//	amqp.AddConsumer[EventT](broker, cfg, handler)
//	broker.Use(amqp.WithRecover(), amqp.WithLogging())  // group middlewares
//	broker.Run(ctx)  // blocks until ctx is done or a consumer fails
//
// Standalone mode: if conn is nil, Publish returns an error
// and Run blocks until ctx is done without starting consumers.
type Broker struct {
	publishers *publisherManager
	consumers  *consumerGroup
}

// NewBroker creates a ready-to-use Broker backed by the given connection.
// Pass nil for standalone mode (no AMQP server required).
func NewBroker(conn *amqp091.Connection, poolCfg PoolConfig) *Broker {
	poolCfg = poolCfg.withDefaults()

	return &Broker{
		publishers: newPublisherManager(conn, poolCfg),
		consumers:  newConsumerGroup(conn),
	}
}

// Publish sends a raw message to the given exchange with the specified routing key.
func (b *Broker) Publish(ctx context.Context, exchange, routingKey string, body []byte, g DeliveryGuarantee) error {
	p, err := b.publishers.get(g)
	if err != nil {
		return err
	}

	return p.Publish(ctx, exchange, routingKey, body)
}

// PublishJSON validates the payload struct, marshals it to JSON, and publishes.
func (b *Broker) PublishJSON(ctx context.Context, exchange, routingKey string, payload any, g DeliveryGuarantee) error {
	if err := validate.StructCtx(ctx, payload); err != nil {
		return fmt.Errorf("amqp: validate: %w", err)
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("amqp: marshal: %w", err)
	}

	return b.Publish(ctx, exchange, routingKey, body, g)
}

// Use appends group-level middlewares that apply to every registered consumer.
// Must be called before AddConsumer so that middlewares are captured at registration time.
func (b *Broker) Use(mws ...Middleware) {
	b.consumers.Use(mws...)
}

// Run declares all queues/bindings, then starts consuming on dedicated channels.
// It blocks until ctx is cancelled, Shutdown is called, or a consumer returns an error.
// On shutdown the delivery loop stops immediately, but in-flight handlers finish.
func (b *Broker) Run(ctx context.Context) error {
	return b.consumers.Run(ctx)
}

// Shutdown gracefully stops consumers (no new deliveries, in-flight handlers finish)
// and closes all publisher channels.
func (b *Broker) Shutdown() {
	b.consumers.Shutdown()
	b.publishers.closeAll()
}
