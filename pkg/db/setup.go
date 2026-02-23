package db

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"time"

	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
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
// *slog.Logger parameter ensures Wire initializes the logger before the database.
func Setup(ctx context.Context, cfg DBConfig, _ *slog.Logger) *bun.DB {
	if cfg.Standalone {
		slog.Warn("standalone mode: skipping database connection")
		return nil
	}

	sqldb := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(cfg.dsn())))

	db := bun.NewDB(sqldb, pgdialect.New())

	ctx, cancel := context.WithTimeout(ctx, cfg.ConnectTimeout)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		panic("failed to connect to database: " + err.Error())
	}

	slog.Info("postgres connected", slog.String("host", cfg.Host), slog.Int("port", cfg.Port), slog.String("db", cfg.Name))

	return db
}
