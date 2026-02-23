package testcontainer

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/go-connections/nat"
	amqp091 "github.com/rabbitmq/amqp091-go"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

// AMQPContainer manages a RabbitMQ test container.
type AMQPContainer struct {
	HostPort string

	container testcontainers.Container
	conn      *amqp091.Connection
}

func (a *AMQPContainer) Start(ctx context.Context) error {
	c, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: testcontainers.ContainerRequest{
			Image:        "rabbitmq:3-alpine",
			ExposedPorts: []string{a.HostPort + ":5672/tcp"},
			HostConfigModifier: func(hc *container.HostConfig) {
				hc.PortBindings = nat.PortMap{
					"5672/tcp": []nat.PortBinding{{HostPort: a.HostPort}},
				}
			},
			WaitingFor: wait.ForLog("Server startup complete").
				WithStartupTimeout(60 * time.Second),
		},
		Started: true,
	})
	if err != nil {
		return fmt.Errorf("start rabbitmq container: %w", err)
	}
	a.container = c

	url := fmt.Sprintf("amqp://guest:guest@localhost:%s/", a.HostPort)
	a.conn, err = amqp091.Dial(url)
	if err != nil {
		return fmt.Errorf("connect to rabbitmq: %w", err)
	}

	return nil
}

func (a *AMQPContainer) Migrate(_ context.Context) error { return nil }

func (a *AMQPContainer) InitFixtures() error { return nil }

func (a *AMQPContainer) LoadFixtures(_ *testing.T, _ string) {}

func (a *AMQPContainer) Clean(_ context.Context) error { return nil }

func (a *AMQPContainer) Close() {
	if a.conn != nil {
		_ = a.conn.Close()
	}
}

func (a *AMQPContainer) Terminate(ctx context.Context) {
	if a.container != nil {
		_ = a.container.Terminate(ctx)
	}
}

// Conn returns the AMQP connection.
func (a *AMQPContainer) Conn() *amqp091.Connection {
	return a.conn
}
