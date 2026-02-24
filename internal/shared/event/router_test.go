//go:build unit

package event

import (
	"context"
	"errors"
	"testing"

	pkgamqp "starter-boilerplate/pkg/amqp"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type testEvent struct {
	Value string `json:"value" validate:"required"`
}

func (testEvent) EventName() string { return "test.event" }

func TestRouter_DispatchToRegisteredHandler(t *testing.T) {
	r := NewRouter()

	var received testEvent
	Route(r, func(_ context.Context, payload testEvent, _ pkgamqp.DeliveryMeta) error {
		received = payload
		return nil
	})

	handler := r.Handler()
	err := handler(context.Background(), []byte(`{"value":"hello"}`), pkgamqp.DeliveryMeta{RoutingKey: "test.event"})

	require.NoError(t, err)
	assert.Equal(t, "hello", received.Value)
}

func TestRouter_HandlerError(t *testing.T) {
	r := NewRouter()

	Route(r, func(_ context.Context, _ testEvent, _ pkgamqp.DeliveryMeta) error {
		return errors.New("handler failed")
	})

	handler := r.Handler()
	err := handler(context.Background(), []byte(`{"value":"x"}`), pkgamqp.DeliveryMeta{RoutingKey: "test.event"})

	assert.EqualError(t, err, "handler failed")
}

func TestRouter_UnmarshalError(t *testing.T) {
	r := NewRouter()

	Route(r, func(_ context.Context, _ testEvent, _ pkgamqp.DeliveryMeta) error {
		t.Fatal("handler should not be called")
		return nil
	})

	handler := r.Handler()
	err := handler(context.Background(), []byte(`{invalid json`), pkgamqp.DeliveryMeta{RoutingKey: "test.event"})

	require.Error(t, err)
	assert.Contains(t, err.Error(), "unmarshal")
}

func TestRouter_ValidationError(t *testing.T) {
	r := NewRouter()

	Route(r, func(_ context.Context, _ testEvent, _ pkgamqp.DeliveryMeta) error {
		t.Fatal("handler should not be called")
		return nil
	})

	handler := r.Handler()
	err := handler(context.Background(), []byte(`{"value":""}`), pkgamqp.DeliveryMeta{RoutingKey: "test.event"})

	require.Error(t, err)
	assert.Contains(t, err.Error(), "validate")
}

func TestRouter_UnhandledRoutingKey_NoFallback(t *testing.T) {
	r := NewRouter()
	handler := r.Handler()

	err := handler(context.Background(), []byte(`{}`), pkgamqp.DeliveryMeta{RoutingKey: "unknown.event"})

	assert.NoError(t, err)
}

func TestRouter_UnhandledRoutingKey_WithFallback(t *testing.T) {
	r := NewRouter()

	var fallbackKey string
	r.Default(func(_ context.Context, _ []byte, meta pkgamqp.DeliveryMeta) error {
		fallbackKey = meta.RoutingKey
		return nil
	})

	handler := r.Handler()
	err := handler(context.Background(), []byte(`{}`), pkgamqp.DeliveryMeta{RoutingKey: "unknown.event"})

	assert.NoError(t, err)
	assert.Equal(t, "unknown.event", fallbackKey)
}

func TestRouter_FallbackError(t *testing.T) {
	r := NewRouter()

	r.Default(func(_ context.Context, _ []byte, _ pkgamqp.DeliveryMeta) error {
		return errors.New("fallback error")
	})

	handler := r.Handler()
	err := handler(context.Background(), []byte(`{}`), pkgamqp.DeliveryMeta{RoutingKey: "unknown.event"})

	assert.EqualError(t, err, "fallback error")
}
