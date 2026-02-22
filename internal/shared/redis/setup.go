package redis

import (
	"context"
	"time"

	goredis "github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

type RedisConfig struct {
	Standalone     bool          `yaml:"standalone"`
	Addr           string        `yaml:"addr" validate:"required_unless=Standalone true"`
	Password       string        `yaml:"password"`
	DB             int           `yaml:"db"`
	ConnectTimeout time.Duration `yaml:"connect_timeout" validate:"required_unless=Standalone true"`
}

// Setup creates a new Redis client.
// *zap.Logger parameter ensures Wire initializes the logger before Redis.
func Setup(ctx context.Context, cfg RedisConfig, _ *zap.Logger) *goredis.Client {
	if cfg.Standalone {
		zap.L().Warn("standalone mode: skipping redis connection")
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

	zap.L().Info("redis connected", zap.String("addr", cfg.Addr))

	return client
}
