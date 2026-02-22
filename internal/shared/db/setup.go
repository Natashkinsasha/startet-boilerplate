package db

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
	"go.uber.org/zap"
)

type DBConfig struct {
	Standalone     bool          `yaml:"standalone"`
	Host           string        `yaml:"host" validate:"required_unless=Standalone true"`
	Port           int           `yaml:"port" validate:"required_unless=Standalone true"`
	Name           string        `yaml:"name" validate:"required_unless=Standalone true"`
	User           string        `yaml:"user" validate:"required_unless=Standalone true"`
	Password       string        `yaml:"password" validate:"required_unless=Standalone true"`
	ConnectTimeout time.Duration `yaml:"connect_timeout" validate:"required_unless=Standalone true"`
}

func (d DBConfig) dsn() string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?sslmode=disable",
		d.User, d.Password, d.Host, d.Port, d.Name,
	)
}

// Setup creates a new bun.DB connection.
// *zap.Logger parameter ensures Wire initializes the logger before the database.
func Setup(ctx context.Context, cfg DBConfig, _ *zap.Logger) *bun.DB {
	if cfg.Standalone {
		zap.L().Warn("standalone mode: skipping database connection")
		return nil
	}

	sqldb := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(cfg.dsn())))

	db := bun.NewDB(sqldb, pgdialect.New())

	ctx, cancel := context.WithTimeout(ctx, cfg.ConnectTimeout)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		panic("failed to connect to database: " + err.Error())
	}

	zap.L().Info("postgres connected", zap.String("host", cfg.Host), zap.Int("port", cfg.Port), zap.String("db", cfg.Name))

	return db
}
