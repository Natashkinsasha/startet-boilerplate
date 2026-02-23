package amqp

import (
	"context"
	"fmt"
	"log/slog"
	"sync"

	amqp091 "github.com/rabbitmq/amqp091-go"
	"golang.org/x/sync/errgroup"
)

// consumerGroup manages a set of AMQP consumers internally.
// Not intended for direct use — access via [Broker].
type consumerGroup struct {
	conn      *amqp091.Connection
	consumers []*consumer
	mws       []Middleware

	mu     sync.Mutex
	cancel context.CancelFunc
}

func newConsumerGroup(conn *amqp091.Connection) *consumerGroup {
	return &consumerGroup{conn: conn}
}

func (g *consumerGroup) Use(mws ...Middleware) {
	g.mws = append(g.mws, mws...)
}

// AddConsumer registers a typed consumer in the broker.
// The handler receives a deserialized payload of type T and [DeliveryMeta].
// Group-level middlewares (set via [Broker].Use) are applied before per-consumer mws.
func AddConsumer[T any](b *Broker, cfg ConsumerConfig, fn func(ctx context.Context, payload T, meta DeliveryMeta) error, mws ...Middleware) {
	g := b.consumers
	allMws := append(g.mws[:len(g.mws):len(g.mws)], mws...)
	g.consumers = append(g.consumers, newConsumer(cfg, typedHandler(fn), allMws...))
}

// AddRawConsumer registers a consumer that receives raw message bytes and [DeliveryMeta].
// Group-level middlewares (set via [Broker].Use) are applied before per-consumer mws.
func AddRawConsumer(b *Broker, cfg ConsumerConfig, fn func(ctx context.Context, body []byte, meta DeliveryMeta) error, mws ...Middleware) {
	g := b.consumers
	allMws := append(g.mws[:len(g.mws):len(g.mws)], mws...)
	g.consumers = append(g.consumers, newConsumer(cfg, rawHandler(fn), allMws...))
}

func (g *consumerGroup) Run(ctx context.Context) error {
	if g.conn == nil {
		slog.Warn("standalone mode: skipping amqp consumers")
		<-ctx.Done()
		return nil
	}

	slog.Info("amqp consumers started", slog.Int("count", len(g.consumers)))

	// Declare queues/bindings once per consumer before starting goroutines.
	for _, c := range g.consumers {
		ch, err := g.conn.Channel()
		if err != nil {
			return fmt.Errorf("open channel for declare: %w", err)
		}
		if err := c.Declare(ch); err != nil {
			ch.Close()
			return fmt.Errorf("declare consumer %q: %w", c.cfg.Queue, err)
		}
		ch.Close()
	}

	eg, egCtx := errgroup.WithContext(ctx)

	// consumeCtx controls message delivery.
	// Cancelled by Shutdown() or when any consumer goroutine fails.
	consumeCtx, consumeCancel := context.WithCancel(egCtx)
	defer consumeCancel()

	g.mu.Lock()
	g.cancel = consumeCancel
	g.mu.Unlock()

	for _, c := range g.consumers {
		for range c.cfg.concurrency() {
			eg.Go(func() error {
				ch, err := g.conn.Channel()
				if err != nil {
					return err
				}
				defer ch.Close()

				return c.Consume(consumeCtx, ch)
			})
		}
	}

	err := eg.Wait()
	slog.Info("amqp consumers stopped")
	return err
}

func (g *consumerGroup) Shutdown() {
	g.mu.Lock()
	defer g.mu.Unlock()

	if g.cancel != nil {
		g.cancel()
	}
}
