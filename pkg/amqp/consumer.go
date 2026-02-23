package amqp

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"reflect"

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
	Concurrency          int    // parallel goroutines per consumer, default: 1
	RetryOnError         *bool  // nack+requeue on error, default: false
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

func (c ConsumerConfig) concurrency() int {
	if c.Concurrency <= 0 {
		return 1
	}
	return c.Concurrency
}

func (c ConsumerConfig) retryOnError() bool {
	if c.RetryOnError == nil {
		return false
	}
	return *c.RetryOnError
}

// consumer handles queue declaration, binding, prefetch, consume loop, and ack/nack
// for a single AMQP queue. Created internally by [AddConsumer].
type consumer struct {
	cfg     ConsumerConfig
	handler HandlerFunc
}

func newConsumer(cfg ConsumerConfig, handler HandlerFunc, mws ...Middleware) *consumer {
	return &consumer{cfg: cfg, handler: Chain(handler, mws...)}
}

// DeliveryMeta holds AMQP message metadata extracted from amqp091.Delivery.
type DeliveryMeta struct {
	Headers     amqp091.Table
	RoutingKey  string
	Exchange    string
	MessageID   string
	ContentType string
	Timestamp   int64
}

func newDeliveryMeta(msg *amqp091.Delivery) DeliveryMeta {
	return DeliveryMeta{
		Headers:     msg.Headers,
		RoutingKey:  msg.RoutingKey,
		Exchange:    msg.Exchange,
		MessageID:   msg.MessageId,
		ContentType: msg.ContentType,
		Timestamp:   msg.Timestamp.Unix(),
	}
}

// typedHandler wraps a typed handler function into a HandlerFunc.
// It unmarshals the message body into T and validates it before calling fn.
// The handler receives the deserialized payload followed by DeliveryMeta.
func typedHandler[T any](fn func(ctx context.Context, payload T, meta DeliveryMeta) error) HandlerFunc {
	return func(ctx context.Context, msg amqp091.Delivery) error {
		var payload T
		if err := json.Unmarshal(msg.Body, &payload); err != nil {
			return fmt.Errorf("unmarshal: %w", err)
		}
		if isStruct(payload) {
			if err := validate.StructCtx(ctx, payload); err != nil {
				return fmt.Errorf("validate: %w", err)
			}
		}
		return fn(ctx, payload, newDeliveryMeta(&msg))
	}
}

func isStruct(v any) bool {
	t := reflect.TypeOf(v)
	if t == nil {
		return false
	}
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	return t.Kind() == reflect.Struct
}

func (c *consumer) Declare(ch *amqp091.Channel) error {
	if err := c.declare(ch); err != nil {
		return err
	}
	return c.bind(ch)
}

func (c *consumer) Consume(ctx context.Context, ch *amqp091.Channel) error {
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

func (c *consumer) declare(ch *amqp091.Channel) error {
	var args amqp091.Table
	if c.cfg.DeadLetterExchange != "" {
		if err := c.declareDLQ(ch); err != nil {
			return fmt.Errorf("declare dlq: %w", err)
		}
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

func (c *consumer) declareDLQ(ch *amqp091.Channel) error {
	if err := ch.ExchangeDeclare(c.cfg.DeadLetterExchange, "topic", true, false, false, false, nil); err != nil {
		return err
	}

	dlqQueue := c.cfg.Queue + ".dlq"
	if _, err := ch.QueueDeclare(dlqQueue, true, false, false, false, nil); err != nil {
		return err
	}

	routingKey := c.cfg.DeadLetterRoutingKey
	if routingKey == "" {
		routingKey = c.cfg.Queue
	}
	return ch.QueueBind(dlqQueue, routingKey, c.cfg.DeadLetterExchange, false, nil)
}

func (c *consumer) bind(ch *amqp091.Channel) error {
	if c.cfg.Exchange == "" {
		return nil
	}
	return ch.QueueBind(c.cfg.Queue, c.cfg.RoutingKey, c.cfg.Exchange, false, nil)
}

func (c *consumer) process(parent context.Context, msg amqp091.Delivery) {
	ctx, cancel := context.WithCancel(context.WithoutCancel(parent))
	defer cancel()

	err := c.handler(ctx, msg)
	if err != nil {
		_ = msg.Nack(false, c.cfg.retryOnError())
		return
	}
	if ackErr := msg.Ack(false); ackErr != nil {
		slog.Error("failed to ack message", slog.Any("error", ackErr))
	}
}
