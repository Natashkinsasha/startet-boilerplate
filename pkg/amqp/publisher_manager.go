package amqp

import (
	"fmt"
	"log/slog"

	amqp091 "github.com/rabbitmq/amqp091-go"
)

// publisherManager owns a dedicated publisher pool for each delivery guarantee.
type publisherManager struct {
	atLeastOnce publisher
	atMostOnce  publisher
}

func newPublisherManager(conn *amqp091.Connection, poolCfg PoolConfig) *publisherManager {
	if conn == nil {
		return &publisherManager{}
	}

	return &publisherManager{
		atLeastOnce: newPublisherPool(conn, AtLeastOnce, poolCfg),
		atMostOnce:  newPublisherPool(conn, AtMostOnce, poolCfg),
	}
}

// get returns the pool for the given delivery guarantee.
func (m *publisherManager) get(g DeliveryGuarantee) (publisher, error) {
	switch g {
	case AtLeastOnce:
		if m.atLeastOnce == nil {
			return nil, fmt.Errorf("amqp: connection is nil (standalone mode?)")
		}
		return m.atLeastOnce, nil
	case AtMostOnce:
		if m.atMostOnce == nil {
			return nil, fmt.Errorf("amqp: connection is nil (standalone mode?)")
		}
		return m.atMostOnce, nil
	default:
		return nil, fmt.Errorf("amqp: unknown delivery guarantee: %d", g)
	}
}

// closeAll closes every publisher pool.
func (m *publisherManager) closeAll() {
	for g, p := range map[DeliveryGuarantee]publisher{
		AtLeastOnce: m.atLeastOnce,
		AtMostOnce:  m.atMostOnce,
	} {
		if p == nil {
			continue
		}
		if err := p.Close(); err != nil {
			slog.Error("failed to close publisher", slog.Any("guarantee", g), slog.Any("error", err))
		}
	}
}
