package amqp

import (
	"context"
	"fmt"
	"log/slog"
	"sync"

	amqp091 "github.com/rabbitmq/amqp091-go"
	"golang.org/x/sync/errgroup"
)

// ConsumerGroup is a self-contained component that manages a set of AMQP
// consumers. It holds the connection internally and exposes Run / Shutdown.
type ConsumerGroup struct {
	conn      *amqp091.Connection
	consumers []*Consumer
	mws       []Middleware

	mu     sync.Mutex
	cancel context.CancelFunc
}

// NewConsumerGroup creates a ConsumerGroup with the given AMQP connection.
// If conn is nil, Run will log a warning and block until ctx is done (standalone mode).
func NewConsumerGroup(conn *amqp091.Connection) *ConsumerGroup {
	return &ConsumerGroup{conn: conn}
}

// Use sets group-level middlewares that apply to all consumers.
// Group middlewares run before per-consumer middlewares.
func (g *ConsumerGroup) Use(mws ...Middleware) {
	g.mws = append(g.mws, mws...)
}

// AddConsumer creates a typed consumer and registers it in the group.
// Group middlewares are prepended before per-consumer middlewares.
func AddConsumer[T any](g *ConsumerGroup, cfg ConsumerConfig, fn func(ctx context.Context, payload T, meta DeliveryMeta) error, mws ...Middleware) {
	allMws := append(g.mws[:len(g.mws):len(g.mws)], mws...)
	g.consumers = append(g.consumers, NewConsumer(cfg, TypedHandler(fn), allMws...))
}

// Merge appends consumers from other groups into this group.
func (g *ConsumerGroup) Merge(others ...*ConsumerGroup) {
	for _, o := range others {
		g.consumers = append(g.consumers, o.consumers...)
	}
}

// Run starts every consumer on a dedicated AMQP channel and blocks until
// ctx is cancelled, Shutdown is called, or any consumer returns an error.
//
// On shutdown the delivery loop stops immediately, but in-flight message
// handlers are allowed to finish before Run returns.
func (g *ConsumerGroup) Run(ctx context.Context) error {
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

// Shutdown gracefully stops all consumers started by Run.
// It stops accepting new messages but lets in-flight handlers finish.
func (g *ConsumerGroup) Shutdown() {
	g.mu.Lock()
	defer g.mu.Unlock()

	if g.cancel != nil {
		g.cancel()
	}
}
