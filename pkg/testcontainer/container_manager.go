package testcontainer

import (
	"context"
	"testing"
)

// ContainerManager is a composite Container that manages multiple children.
type ContainerManager struct {
	containers []Container
}

// NewContainerManager creates a composite from any number of containers.
func NewContainerManager(containers ...Container) *ContainerManager {
	return &ContainerManager{containers: containers}
}

func (m *ContainerManager) Start(ctx context.Context) error {
	for _, c := range m.containers {
		if err := c.Start(ctx); err != nil {
			return err
		}
	}
	return nil
}

func (m *ContainerManager) Migrate(ctx context.Context) error {
	for _, c := range m.containers {
		if err := c.Migrate(ctx); err != nil {
			return err
		}
	}
	return nil
}

func (m *ContainerManager) InitFixtures() error {
	for _, c := range m.containers {
		if err := c.InitFixtures(); err != nil {
			return err
		}
	}
	return nil
}

func (m *ContainerManager) LoadFixtures(t *testing.T, key string) {
	t.Helper()
	for _, c := range m.containers {
		c.LoadFixtures(t, key)
	}
}

func (m *ContainerManager) Clean(ctx context.Context) error {
	for _, c := range m.containers {
		if err := c.Clean(ctx); err != nil {
			return err
		}
	}
	return nil
}

func (m *ContainerManager) Close() {
	for _, c := range m.containers {
		c.Close()
	}
}

func (m *ContainerManager) Terminate(ctx context.Context) {
	for _, c := range m.containers {
		c.Terminate(ctx)
	}
}
