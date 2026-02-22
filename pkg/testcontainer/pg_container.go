package testcontainer

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"testing"
	"time"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/go-connections/nat"
	"github.com/go-testfixtures/testfixtures/v3"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
)

// PgContainer manages a postgres test container with bun.DB, migrations, and fixtures.
type PgContainer struct {
	Database string
	Username string
	Password string
	HostPort string
	Fixtures map[string]FixtureSet

	container *postgres.PostgresContainer
	db        *bun.DB
	loaders   map[string]*testfixtures.Loader
}

func (p *PgContainer) Start(ctx context.Context) error {
	c, err := postgres.Run(ctx, "postgres:16-alpine",
		postgres.WithDatabase(p.Database),
		postgres.WithUsername(p.Username),
		postgres.WithPassword(p.Password),
		testcontainers.CustomizeRequest(testcontainers.GenericContainerRequest{
			ContainerRequest: testcontainers.ContainerRequest{
				HostConfigModifier: func(hc *container.HostConfig) {
					hc.PortBindings = nat.PortMap{
						"5432/tcp": []nat.PortBinding{{HostPort: p.HostPort}},
					}
				},
			},
		}),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(30*time.Second),
		),
	)
	if err != nil {
		return fmt.Errorf("start postgres container: %w", err)
	}
	p.container = c

	dsn := fmt.Sprintf("postgres://%s:%s@localhost:%s/%s?sslmode=disable",
		p.Username, p.Password, p.HostPort, p.Database)
	sqldb := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(dsn)))
	p.db = bun.NewDB(sqldb, pgdialect.New())

	return nil
}

func (p *PgContainer) Migrate(ctx context.Context) error {
	migrationsDir := "migrations"

	entries, err := os.ReadDir(migrationsDir)
	if err != nil {
		return fmt.Errorf("read migrations dir: %w", err)
	}

	var upFiles []string
	for _, e := range entries {
		if !e.IsDir() && strings.HasSuffix(e.Name(), ".up.sql") {
			upFiles = append(upFiles, e.Name())
		}
	}
	sort.Strings(upFiles)

	for _, name := range upFiles {
		data, err := os.ReadFile(filepath.Join(migrationsDir, name))
		if err != nil {
			return fmt.Errorf("read migration %s: %w", name, err)
		}
		if _, err := p.db.ExecContext(ctx, string(data)); err != nil {
			return fmt.Errorf("exec migration %s: %w", name, err)
		}
	}

	return nil
}

func (p *PgContainer) InitFixtures() error {
	p.loaders = make(map[string]*testfixtures.Loader, len(p.Fixtures))

	for key, fs := range p.Fixtures {
		opts := []func(*testfixtures.Loader) error{
			testfixtures.Database(p.db.DB),
			testfixtures.Dialect("postgresql"),
			testfixtures.Template(),
		}
		if fs.TemplateData != nil {
			opts = append(opts, testfixtures.TemplateData(fs.TemplateData))
		}
		opts = append(opts, testfixtures.Directory(fs.Dir))

		loader, err := testfixtures.New(opts...)
		if err != nil {
			return fmt.Errorf("init fixtures %q: %w", key, err)
		}
		p.loaders[key] = loader
	}

	return nil
}

func (p *PgContainer) LoadFixtures(t *testing.T, key string) {
	t.Helper()
	loader, ok := p.loaders[key]
	require.True(t, ok, "fixture set %q not found", key)
	err := loader.Load()
	require.NoError(t, err, "failed to load fixtures %q", key)
}

func (p *PgContainer) Clean(ctx context.Context) error {
	_, err := p.db.ExecContext(ctx, `
		DO $$ DECLARE t text;
		BEGIN
			FOR t IN
				SELECT tablename FROM pg_tables WHERE schemaname = 'public'
			LOOP
				EXECUTE 'TRUNCATE TABLE public.' || quote_ident(t) || ' CASCADE';
			END LOOP;
		END $$;
	`)
	return err
}

func (p *PgContainer) Close() {
	if p.db != nil {
		_ = p.db.Close()
	}
}

func (p *PgContainer) Terminate(ctx context.Context) {
	if p.container != nil {
		_ = p.container.Terminate(ctx)
	}
}

// DB returns the bun.DB instance.
func (p *PgContainer) DB() *bun.DB {
	return p.db
}
