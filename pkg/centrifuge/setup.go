package centrifuge

import (
	"context"
	"log/slog"
	"time"

	"github.com/centrifugal/centrifuge"
	goredis "github.com/redis/go-redis/v9"
)

type Config struct {
	Standalone  bool          `yaml:"standalone"`
	HistorySize int           `yaml:"history_size"`
	HistoryTTL  time.Duration `yaml:"history_ttl"`
}

// Setup creates a centrifuge Node backed by Redis for horizontal scaling.
// Returns nil in standalone mode.
// *slog.Logger parameter ensures Wire initializes the logger before centrifuge.
func Setup(_ context.Context, cfg Config, redisClient *goredis.Client, _ *slog.Logger) *centrifuge.Node {
	if cfg.Standalone {
		slog.Warn("standalone mode: skipping centrifuge node")
		return nil
	}

	node, err := centrifuge.New(centrifuge.Config{
		LogLevel: centrifuge.LogLevelInfo,
		LogHandler: func(e centrifuge.LogEntry) {
			slog.Log(context.Background(), toSlogLevel(e.Level), e.Message,
				slog.Any("fields", e.Fields),
			)
		},
	})
	if err != nil {
		panic("centrifuge: create node: " + err.Error())
	}

	opts := redisClient.Options()
	shardCfg := centrifuge.RedisShardConfig{
		Address:  opts.Addr,
		Password: opts.Password,
		DB:       opts.DB,
	}

	shard, err := centrifuge.NewRedisShard(node, shardCfg)
	if err != nil {
		panic("centrifuge: create redis shard: " + err.Error())
	}

	broker, err := centrifuge.NewRedisBroker(node, centrifuge.RedisBrokerConfig{
		Shards: []*centrifuge.RedisShard{shard},
	})
	if err != nil {
		panic("centrifuge: create redis broker: " + err.Error())
	}
	node.SetBroker(broker)

	presenceMgr, err := centrifuge.NewRedisPresenceManager(node, centrifuge.RedisPresenceManagerConfig{
		Shards: []*centrifuge.RedisShard{shard},
	})
	if err != nil {
		panic("centrifuge: create redis presence manager: " + err.Error())
	}
	node.SetPresenceManager(presenceMgr)

	slog.Info("centrifuge node created", slog.String("redis_addr", opts.Addr))
	return node
}

func toSlogLevel(lvl centrifuge.LogLevel) slog.Level {
	switch lvl {
	case centrifuge.LogLevelTrace, centrifuge.LogLevelDebug:
		return slog.LevelDebug
	case centrifuge.LogLevelInfo:
		return slog.LevelInfo
	case centrifuge.LogLevelWarn:
		return slog.LevelWarn
	case centrifuge.LogLevelError:
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}
