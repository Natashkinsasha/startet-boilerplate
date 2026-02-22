package testcontainer

import "context"

// Setup starts all containers, runs migrations, and initialises fixtures.
func Setup(ctx context.Context, containers ...Container) (*ContainerManager, error) {
	m := NewContainerManager(containers...)

	if err := m.Start(ctx); err != nil {
		return nil, err
	}
	if err := m.Migrate(ctx); err != nil {
		return m, err
	}
	if err := m.InitFixtures(); err != nil {
		return m, err
	}

	return m, nil
}

// SetupPgContainer creates, starts, and migrates a single PgContainer.
func SetupPgContainer(ctx context.Context, pg *PgContainer) (*PgContainer, error) {
	if err := pg.Start(ctx); err != nil {
		return nil, err
	}
	if err := pg.Migrate(ctx); err != nil {
		return pg, err
	}
	if err := pg.InitFixtures(); err != nil {
		return pg, err
	}
	return pg, nil
}
