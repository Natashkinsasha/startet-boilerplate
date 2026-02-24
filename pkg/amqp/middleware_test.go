//go:build unit

package amqp

import (
	"context"
	"testing"

	amqp091 "github.com/rabbitmq/amqp091-go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestChain_Empty(t *testing.T) {
	called := false
	handler := func(_ context.Context, _ amqp091.Delivery) error {
		called = true
		return nil
	}

	chained := Chain(handler)

	err := chained(context.Background(), amqp091.Delivery{})
	require.NoError(t, err)
	assert.True(t, called)
}

func TestChain_Order(t *testing.T) {
	var order []string

	mw := func(name string) Middleware {
		return func(next HandlerFunc) HandlerFunc {
			return func(ctx context.Context, msg amqp091.Delivery) error {
				order = append(order, name+":before")
				err := next(ctx, msg)
				order = append(order, name+":after")
				return err
			}
		}
	}

	handler := func(_ context.Context, _ amqp091.Delivery) error {
		order = append(order, "handler")
		return nil
	}

	chained := Chain(handler, mw("first"), mw("second"))

	err := chained(context.Background(), amqp091.Delivery{})
	require.NoError(t, err)
	assert.Equal(t, []string{
		"first:before",
		"second:before",
		"handler",
		"second:after",
		"first:after",
	}, order)
}

func TestChain_SingleMiddleware(t *testing.T) {
	var trace []string

	mw := func(next HandlerFunc) HandlerFunc {
		return func(ctx context.Context, msg amqp091.Delivery) error {
			trace = append(trace, "mw")
			return next(ctx, msg)
		}
	}

	handler := func(_ context.Context, _ amqp091.Delivery) error {
		trace = append(trace, "handler")
		return nil
	}

	chained := Chain(handler, mw)

	err := chained(context.Background(), amqp091.Delivery{})
	require.NoError(t, err)
	assert.Equal(t, []string{"mw", "handler"}, trace)
}
