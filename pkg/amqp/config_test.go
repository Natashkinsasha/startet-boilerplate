//go:build unit

package amqp

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPoolConfig_WithDefaults_ZeroValues(t *testing.T) {
	cfg := PoolConfig{}.withDefaults()

	assert.Equal(t, 1, cfg.InitialSize)
	assert.Equal(t, 5, cfg.MaxSize)
}

func TestPoolConfig_WithDefaults_NegativeValues(t *testing.T) {
	cfg := PoolConfig{InitialSize: -1, MaxSize: -3}.withDefaults()

	assert.Equal(t, 1, cfg.InitialSize)
	assert.Equal(t, 5, cfg.MaxSize)
}

func TestPoolConfig_WithDefaults_CustomValues(t *testing.T) {
	cfg := PoolConfig{InitialSize: 3, MaxSize: 10}.withDefaults()

	assert.Equal(t, 3, cfg.InitialSize)
	assert.Equal(t, 10, cfg.MaxSize)
}

func TestPoolConfig_WithDefaults_MaxLessThanInitial(t *testing.T) {
	cfg := PoolConfig{InitialSize: 8, MaxSize: 2}.withDefaults()

	assert.Equal(t, 8, cfg.InitialSize)
	assert.Equal(t, 8, cfg.MaxSize)
}
