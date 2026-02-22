package testcontainer

import (
	"context"
	"testing"
)

// FixtureSet describes a named set of test fixtures.
type FixtureSet struct {
	Dir          string
	TemplateData map[string]interface{}
}

// Container defines the lifecycle of a test infrastructure container.
type Container interface {
	Start(ctx context.Context) error
	Migrate(ctx context.Context) error
	InitFixtures() error
	LoadFixtures(t *testing.T, key string)
	Clean(ctx context.Context) error
	Close()
	Terminate(ctx context.Context)
}
