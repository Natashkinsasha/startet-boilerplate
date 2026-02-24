package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"strings"

	"starter-boilerplate/internal/shared/config"
	"starter-boilerplate/internal/shared/logger"
	"starter-boilerplate/migrations"
	pkgdb "starter-boilerplate/pkg/db"
	pkgmigrate "starter-boilerplate/pkg/migrate"

	"github.com/uptrace/bun"
)

func main() {
	cfg := config.SetupConfig()
	logger.SetupLogger(cfg.Logger)

	ctx := context.Background()
	db := setupDB(ctx, cfg)
	defer db.Close()

	runner := pkgmigrate.NewRunner(db, migrations.Migrations)

	cmd, arg := parseArgs()
	if err := run(ctx, runner, cmd, arg); err != nil {
		slog.Error("migration failed", slog.Any("error", err))
		os.Exit(1)
	}
}

func run(ctx context.Context, runner *pkgmigrate.Runner, cmd, arg string) error {
	switch cmd {
	case "init":
		return runner.Init(ctx)
	case "migrate":
		return runner.Migrate(ctx)
	case "rollback":
		return runner.Rollback(ctx)
	case "status":
		return runner.Status(ctx)
	case "create":
		if arg == "" {
			return fmt.Errorf("migration name is required: create <name>")
		}
		return runner.CreateSQL(ctx, arg)
	default:
		printUsage()
		return fmt.Errorf("unknown command: %s", cmd)
	}
}

func parseArgs() (cmd, arg string) {
	args := os.Args[1:]
	if len(args) == 0 {
		printUsage()
		os.Exit(1)
	}
	cmd = args[0]
	if len(args) > 1 {
		arg = strings.Join(args[1:], "_")
	}
	return
}

func printUsage() {
	fmt.Println("Usage: migrate <command>")
	fmt.Println()
	fmt.Println("Commands:")
	fmt.Println("  init       Create migration tables")
	fmt.Println("  migrate    Run pending migrations")
	fmt.Println("  rollback   Rollback last migration group")
	fmt.Println("  status     Show migration status")
	fmt.Println("  create     Create new SQL migration files (create <name>)")
}

func setupDB(ctx context.Context, cfg *config.Config) *bun.DB {
	db := pkgdb.Setup(ctx, cfg.DB, slog.Default())
	if db == nil {
		slog.Error("database connection is required for migrations (check db config, standalone must be false)")
		os.Exit(1)
	}
	return db
}
