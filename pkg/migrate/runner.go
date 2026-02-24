package migrate

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/uptrace/bun"
	"github.com/uptrace/bun/migrate"
)

type Runner struct {
	migrator *migrate.Migrator
}

func NewRunner(db *bun.DB, migrations *migrate.Migrations) *Runner {
	return &Runner{
		migrator: migrate.NewMigrator(db, migrations),
	}
}

func (r *Runner) Init(ctx context.Context) error {
	slog.Info("creating migration tables")
	return r.migrator.Init(ctx)
}

func (r *Runner) Migrate(ctx context.Context) error {
	group, err := r.migrator.Migrate(ctx)
	if err != nil {
		return err
	}
	if group.IsZero() {
		slog.Info("no new migrations to run (database is up to date)")
		return nil
	}
	slog.Info("migrated", slog.String("group", group.String()))
	return nil
}

func (r *Runner) Rollback(ctx context.Context) error {
	group, err := r.migrator.Rollback(ctx)
	if err != nil {
		return err
	}
	if group.IsZero() {
		slog.Info("no migrations to rollback")
		return nil
	}
	slog.Info("rolled back", slog.String("group", group.String()))
	return nil
}

func (r *Runner) Status(ctx context.Context) error {
	ms, err := r.migrator.MigrationsWithStatus(ctx)
	if err != nil {
		return err
	}
	slog.Info("migration status",
		slog.String("unapplied", ms.Unapplied().String()),
		slog.String("last_group", ms.LastGroup().String()),
	)
	return nil
}

func (r *Runner) CreateSQL(ctx context.Context, name string) error {
	files, err := r.migrator.CreateSQLMigrations(ctx, name)
	if err != nil {
		return err
	}
	for _, f := range files {
		fmt.Printf("created migration %s (%s)\n", f.Name, f.Path)
	}
	return nil
}
