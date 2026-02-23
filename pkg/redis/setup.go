package redis

import (
	"context"
	"log/slog"
	"time"

	goredis "github.com/redis/go-redis/v9"
)

type RedisConfig struct {
	Standalone     bool          `yaml:"standalone"`
	Addr           string        `yaml:"addr" validate:"required_unless=Standalone true"`
	Password       string        `yaml:"password"`
	DB             int           `yaml:"db"`
	ConnectTimeout time.Duration `yaml:"connect_timeout" validate:"required_unless=Standalone true"`
}

// Setup creates a new Redis client.
// *slog.Logger parameter ensures Wire initializes the logger before Redis.
func Setup(ctx context.Context, cfg RedisConfig, _ *slog.Logger) *goredis.Client {
	if cfg.Standalone {
		slog.Warn("standalone mode: skipping redis connection")
		return nil
	}

	client := goredis.NewClient(&goredis.Options{
		Addr:     cfg.Addr,
		Password: cfg.Password,
		DB:       cfg.DB,
	})

	ctx, cancel := context.WithTimeout(ctx, cfg.ConnectTimeout)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		panic("failed to connect to redis: " + err.Error())
	}

	slog.Info("redis connected", slog.String("addr", cfg.Addr))

	return client
}
