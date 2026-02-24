package outbox

import (
	"context"
	"log/slog"
	"time"

	pkgdb "starter-boilerplate/pkg/db"

	"github.com/uptrace/bun"
)

// Publisher publishes an outbox entry to the message broker.
// The implementation decides the transport (AMQP, Kafka, etc).
type Publisher interface {
	Publish(ctx context.Context, entry Entry) error
}

// RelayConfig controls the relay polling behaviour.
type RelayConfig struct {
	PollInterval time.Duration `yaml:"poll_interval"`
	BatchSize    int           `yaml:"batch_size"`
}

func (c RelayConfig) withDefaults() RelayConfig {
	if c.PollInterval == 0 {
		c.PollInterval = time.Second
	}
	if c.BatchSize == 0 {
		c.BatchSize = 100
	}
	return c
}

// Relay polls the outbox table and publishes entries via Publisher.
type Relay struct {
	db        *bun.DB
	repo      *Repository
	publisher Publisher
	cfg       RelayConfig
}

func NewRelay(db *bun.DB, repo *Repository, publisher Publisher, cfg RelayConfig) *Relay {
	return &Relay{
		db:        db,
		repo:      repo,
		publisher: publisher,
		cfg:       cfg.withDefaults(),
	}
}

// Run polls the outbox on a ticker until ctx is cancelled.
func (r *Relay) Run(ctx context.Context) error {
	slog.Info("outbox relay started",
		slog.Duration("poll_interval", r.cfg.PollInterval),
		slog.Int("batch_size", r.cfg.BatchSize),
	)

	ticker := time.NewTicker(r.cfg.PollInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return nil
		case <-ticker.C:
			if err := r.poll(ctx); err != nil {
				slog.Error("outbox relay poll failed", slog.String("error", err.Error()))
			}
		}
	}
}

func (r *Relay) poll(ctx context.Context) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback() }()

	txCtx := pkgdb.WithTx(ctx, tx)

	entries, err := r.repo.FetchUnpublished(txCtx, r.cfg.BatchSize)
	if err != nil {
		return err
	}
	if len(entries) == 0 {
		return nil
	}

	var published []int64
	for i := range entries {
		if err := r.publisher.Publish(ctx, entries[i]); err != nil {
			slog.Error("outbox relay publish failed",
				slog.Int64("entry_id", entries[i].ID),
				slog.String("error", err.Error()),
			)
			break // stop at first failure to preserve FIFO ordering
		}
		published = append(published, entries[i].ID)
	}

	if len(published) == 0 {
		return nil
	}

	if err := r.repo.MarkPublished(txCtx, published); err != nil {
		return err
	}

	return tx.Commit()
}
