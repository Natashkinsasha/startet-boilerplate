package amqp

import (
	"context"
	"fmt"
	"sync/atomic"

	amqp091 "github.com/rabbitmq/amqp091-go"
)

// publisherPool implements the publisher interface by managing a pool of
// single-channel publishers. It auto-scales from an initial size up to maxSize.
type publisherPool struct {
	conn    *amqp091.Connection
	g       DeliveryGuarantee
	slots   chan publisher // buffered to maxSize
	size    atomic.Int32   // current count of created publishers
	maxSize int32
}

func newPublisherPool(conn *amqp091.Connection, g DeliveryGuarantee, cfg PoolConfig) *publisherPool {
	cfg = cfg.withDefaults()

	p := &publisherPool{
		conn:    conn,
		g:       g,
		slots:   make(chan publisher, cfg.MaxSize),
		maxSize: int32(cfg.MaxSize),
	}

	for range cfg.InitialSize {
		p.slots <- newSinglePublisher(conn, g)
		p.size.Add(1)
	}

	return p
}

// get acquires a publisher from the pool using three-tier logic:
//  1. Non-blocking read from slots (fast path).
//  2. If under maxSize, create a new publisher (grow path).
//  3. Blocking wait on slots with context cancellation (wait path).
func (p *publisherPool) get(ctx context.Context) (publisher, error) {
	// Fast path: grab an idle publisher.
	select {
	case pub := <-p.slots:
		return pub, nil
	default:
	}

	// Grow path: try to create a new publisher if under max.
	for {
		cur := p.size.Load()
		if cur >= p.maxSize {
			break
		}
		if p.size.CompareAndSwap(cur, cur+1) {
			return newSinglePublisher(p.conn, p.g), nil
		}
	}

	// Wait path: block until one is returned or context is done.
	select {
	case pub := <-p.slots:
		return pub, nil
	case <-ctx.Done():
		return nil, fmt.Errorf("amqp: pool acquire cancelled: %w", ctx.Err())
	}
}

func (p *publisherPool) put(pub publisher) {
	p.slots <- pub
}

func (p *publisherPool) Publish(ctx context.Context, exchange, routingKey string, headers amqp091.Table, body []byte) error {
	pub, err := p.get(ctx)
	if err != nil {
		return err
	}
	defer p.put(pub)

	return pub.Publish(ctx, exchange, routingKey, headers, body)
}

func (p *publisherPool) Close() error {
	// Drain all idle publishers.
	for {
		select {
		case pub := <-p.slots:
			_ = pub.Close()
		default:
			return nil
		}
	}
}
