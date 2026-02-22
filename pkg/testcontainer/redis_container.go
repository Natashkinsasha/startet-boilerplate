package testcontainer

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/go-connections/nat"
	goredis "github.com/redis/go-redis/v9"
	"github.com/testcontainers/testcontainers-go"
	tcredis "github.com/testcontainers/testcontainers-go/modules/redis"
	"github.com/testcontainers/testcontainers-go/wait"
)

// RedisContainer manages a redis test container with a go-redis client.
type RedisContainer struct {
	HostPort string

	container *tcredis.RedisContainer
	client    *goredis.Client
}

func (r *RedisContainer) Start(ctx context.Context) error {
	c, err := tcredis.Run(ctx, "redis:7-alpine",
		testcontainers.CustomizeRequest(testcontainers.GenericContainerRequest{
			ContainerRequest: testcontainers.ContainerRequest{
				HostConfigModifier: func(hc *container.HostConfig) {
					hc.PortBindings = nat.PortMap{
						"6379/tcp": []nat.PortBinding{{HostPort: r.HostPort}},
					}
				},
			},
		}),
		testcontainers.WithWaitStrategy(
			wait.ForLog("Ready to accept connections").
				WithStartupTimeout(30*time.Second),
		),
	)
	if err != nil {
		return fmt.Errorf("start redis container: %w", err)
	}
	r.container = c

	r.client = goredis.NewClient(&goredis.Options{
		Addr: "localhost:" + r.HostPort,
	})

	return nil
}

func (r *RedisContainer) Migrate(_ context.Context) error { return nil }

func (r *RedisContainer) InitFixtures() error { return nil }

func (r *RedisContainer) LoadFixtures(_ *testing.T, _ string) {}

func (r *RedisContainer) Clean(ctx context.Context) error {
	return r.client.FlushAll(ctx).Err()
}

func (r *RedisContainer) Close() {
	if r.client != nil {
		_ = r.client.Close()
	}
}

func (r *RedisContainer) Terminate(ctx context.Context) {
	if r.container != nil {
		_ = r.container.Terminate(ctx)
	}
}

// Client returns the go-redis client.
func (r *RedisContainer) Client() *goredis.Client {
	return r.client
}
