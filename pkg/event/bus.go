package event

import "context"

// Event is a domain event that knows its own routing key.
type Event interface {
	EventName() string
}

// Bus publishes domain events. Implementations handle serialization and transport.
type Bus interface {
	Publish(ctx context.Context, event Event) error
}
