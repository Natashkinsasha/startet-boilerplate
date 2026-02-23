package event

import "context"

// Event is a domain event that knows its own routing key.
type Event interface {
	EventName() string
}

// Taggable is an optional interface that events can implement to declare tags.
// Tags are used for server-side filtering via RabbitMQ headers exchange.
type Taggable interface {
	Tags() []string
}

// Bus publishes domain events. Implementations handle serialization and transport.
type Bus interface {
	Publish(ctx context.Context, event Event) error
}
