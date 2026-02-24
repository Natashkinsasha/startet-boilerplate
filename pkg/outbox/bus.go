package outbox

import (
	"context"
	"encoding/json"
	"fmt"
	"time"
)

// Event is a domain event that knows its own routing key.
// Mirrors pkg/event.Event — Go structural typing keeps them compatible.
type Event interface {
	EventName() string
}

// Taggable is an optional interface for events that declare tags.
type Taggable interface {
	Tags() []string
}

// Bus publishes domain events.
type Bus interface {
	Publish(ctx context.Context, event Event) error
}

// OutboxBus implements Bus by inserting events into the outbox table.
// It relies on the transaction being present in context (via pkgdb.WithTx).
type OutboxBus struct {
	outboxRepo *Repository
}

// NewOutboxBus creates an OutboxBus.
func NewOutboxBus(outboxRepo *Repository) *OutboxBus {
	return &OutboxBus{outboxRepo: outboxRepo}
}

func (b *OutboxBus) Publish(ctx context.Context, e Event) error {
	payload, err := json.Marshal(e)
	if err != nil {
		return fmt.Errorf("outbox: marshal event: %w", err)
	}

	headers := make(map[string]any)
	if t, ok := e.(Taggable); ok {
		for _, tag := range t.Tags() {
			headers["tag."+tag] = true
		}
	}

	entry := &Entry{
		EventName: e.EventName(),
		Payload:   payload,
		Headers:   headers,
		CreatedAt: time.Now().Unix(),
	}

	return b.outboxRepo.Insert(ctx, entry)
}
