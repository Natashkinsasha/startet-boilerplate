package amqp

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"runtime/debug"

	"github.com/go-playground/validator/v10"
	amqp091 "github.com/rabbitmq/amqp091-go"
)

var validate = validator.New()

// HandlerFunc processes a single AMQP delivery.
// Return nil to ack, return error to nack.
type HandlerFunc func(ctx context.Context, msg amqp091.Delivery) error

// ConsumerConfig controls queue declaration, binding, and consume behaviour.
type ConsumerConfig struct {
	Queue                string // required
	Exchange             string // optional, for binding
	RoutingKey           string // optional
	Durable              *bool  // default: true
	AutoDelete           bool   // default: false
	Exclusive            bool   // default: false
	PrefetchCount        int    // default: 1
	RetryOnError         *bool  // nack+requeue on error, default: true
	DeadLetterExchange   string // optional, x-dead-letter-exchange
	DeadLetterRoutingKey string // optional, x-dead-letter-routing-key
}

func (c ConsumerConfig) durable() bool {
	if c.Durable == nil {
		return true
	}
	return *c.Durable
}

func (c ConsumerConfig) prefetchCount() int {
	if c.PrefetchCount <= 0 {
		return 1
	}
	return c.PrefetchCount
}

func (c ConsumerConfig) retryOnError() bool {
	if c.RetryOnError == nil {
		return true
	}
	return *c.RetryOnError
}

// Consumer is a reusable AMQP consumer that handles queue declaration,
// binding, prefetch, consume loop, panic recovery, and ack/nack.
type Consumer struct {
	cfg     ConsumerConfig
	handler HandlerFunc
}

func NewConsumer(cfg ConsumerConfig, handler HandlerFunc) *Consumer {
	return &Consumer{cfg: cfg, handler: handler}
}

// TypedHandler wraps a typed handler function into a HandlerFunc.
// It unmarshals the message body into T and validates it before calling fn.
func TypedHandler[T any](fn func(ctx context.Context, payload T) error) HandlerFunc {
	return func(ctx context.Context, msg amqp091.Delivery) error {
		var payload T
		if err := json.Unmarshal(msg.Body, &payload); err != nil {
			return fmt.Errorf("unmarshal: %w", err)
		}
		if err := validate.StructCtx(ctx, payload); err != nil {
			return fmt.Errorf("validate: %w", err)
		}
		return fn(ctx, payload)
	}
}

// Run declares the queue, binds it (if configured), sets QoS, and
// blocks consuming messages until ctx is cancelled or the channel closes.
func (c *Consumer) Run(ctx context.Context, ch *amqp091.Channel) error {
	if err := c.declare(ch); err != nil {
		return err
	}

	if err := c.bind(ch); err != nil {
		return err
	}

	if err := ch.Qos(c.cfg.prefetchCount(), 0, false); err != nil {
		return err
	}

	msgs, err := ch.ConsumeWithContext(ctx, c.cfg.Queue, "", false, c.cfg.Exclusive, false, false, nil)
	if err != nil {
		return err
	}

	for msg := range msgs {
		c.process(ctx, msg)
	}

	return nil
}

func (c *Consumer) declare(ch *amqp091.Channel) error {
	var args amqp091.Table
	if c.cfg.DeadLetterExchange != "" {
		args = amqp091.Table{"x-dead-letter-exchange": c.cfg.DeadLetterExchange}
		if c.cfg.DeadLetterRoutingKey != "" {
			args["x-dead-letter-routing-key"] = c.cfg.DeadLetterRoutingKey
		}
	}

	_, err := ch.QueueDeclare(
		c.cfg.Queue,
		c.cfg.durable(),
		c.cfg.AutoDelete,
		c.cfg.Exclusive,
		false, // noWait
		args,
	)
	return err
}

func (c *Consumer) bind(ch *amqp091.Channel) error {
	if c.cfg.Exchange == "" {
		return nil
	}
	return ch.QueueBind(c.cfg.Queue, c.cfg.RoutingKey, c.cfg.Exchange, false, nil)
}

func (c *Consumer) process(parent context.Context, msg amqp091.Delivery) {
	ctx, cancel := context.WithCancel(parent)
	defer cancel()

	log := slog.Default().With(slog.String("queue", c.cfg.Queue))

	panicked := true
	defer func() {
		if !panicked {
			return
		}
		r := recover()
		log.Error("panic while handling message",
			slog.Any("panic", r),
			slog.String("stack", string(debug.Stack())),
		)
		_ = msg.Nack(false, true)
	}()

	err := c.handler(ctx, msg)
	panicked = false

	if err != nil {
		log.Error("handler error", slog.Any("error", err))
		_ = msg.Nack(false, c.cfg.retryOnError())
		return
	}

	if ackErr := msg.Ack(false); ackErr != nil {
		log.Error("failed to ack message", slog.Any("error", ackErr))
	}
}
