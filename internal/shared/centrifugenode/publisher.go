package centrifugenode

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	pkgcentrifuge "starter-boilerplate/pkg/centrifuge"

	gocentrifuge "github.com/centrifugal/centrifuge"
)

type Publisher struct {
	node        *gocentrifuge.Node
	historySize int
	historyTTL  time.Duration
}

func NewPublisher(node *gocentrifuge.Node, cfg pkgcentrifuge.Config) *Publisher {
	return &Publisher{
		node:        node,
		historySize: cfg.HistorySize,
		historyTTL:  cfg.HistoryTTL,
	}
}

type wsMessage struct {
	Event   string          `json:"event"`
	Payload json.RawMessage `json:"payload"`
}

// Enabled reports whether the centrifuge node is configured.
func (p *Publisher) Enabled() bool {
	return p.node != nil
}

// PublishPersonal publishes an event to the user's personal channel.
// No-op if node is nil (standalone mode).
func (p *Publisher) PublishPersonal(ctx context.Context, userID, eventName string, payload []byte) error {
	if p.node == nil {
		return nil
	}

	data, err := json.Marshal(wsMessage{
		Event:   eventName,
		Payload: payload,
	})
	if err != nil {
		return fmt.Errorf("centrifuge publish: marshal: %w", err)
	}

	var opts []gocentrifuge.PublishOption
	if p.historySize > 0 && p.historyTTL > 0 {
		opts = append(opts, gocentrifuge.WithHistory(p.historySize, p.historyTTL))
	}

	_, err = p.node.Publish("personal:"+userID, data, opts...)
	return err
}
