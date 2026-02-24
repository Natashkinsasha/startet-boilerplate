//go:build unit

package amqp

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func boolPtr(b bool) *bool { return &b }

func TestConsumerConfig_Durable_DefaultTrue(t *testing.T) {
	cfg := ConsumerConfig{}
	assert.True(t, cfg.durable())
}

func TestConsumerConfig_Durable_ExplicitFalse(t *testing.T) {
	cfg := ConsumerConfig{Durable: boolPtr(false)}
	assert.False(t, cfg.durable())
}

func TestConsumerConfig_Durable_ExplicitTrue(t *testing.T) {
	cfg := ConsumerConfig{Durable: boolPtr(true)}
	assert.True(t, cfg.durable())
}

func TestConsumerConfig_PrefetchCount_Default(t *testing.T) {
	cfg := ConsumerConfig{}
	assert.Equal(t, 1, cfg.prefetchCount())
}

func TestConsumerConfig_PrefetchCount_Custom(t *testing.T) {
	cfg := ConsumerConfig{PrefetchCount: 10}
	assert.Equal(t, 10, cfg.prefetchCount())
}

func TestConsumerConfig_PrefetchCount_Negative(t *testing.T) {
	cfg := ConsumerConfig{PrefetchCount: -5}
	assert.Equal(t, 1, cfg.prefetchCount())
}

func TestConsumerConfig_Concurrency_Default(t *testing.T) {
	cfg := ConsumerConfig{}
	assert.Equal(t, 1, cfg.concurrency())
}

func TestConsumerConfig_Concurrency_Custom(t *testing.T) {
	cfg := ConsumerConfig{Concurrency: 4}
	assert.Equal(t, 4, cfg.concurrency())
}

func TestConsumerConfig_RetryOnError_DefaultFalse(t *testing.T) {
	cfg := ConsumerConfig{}
	assert.False(t, cfg.retryOnError())
}

func TestConsumerConfig_RetryOnError_ExplicitTrue(t *testing.T) {
	cfg := ConsumerConfig{RetryOnError: boolPtr(true)}
	assert.True(t, cfg.retryOnError())
}
