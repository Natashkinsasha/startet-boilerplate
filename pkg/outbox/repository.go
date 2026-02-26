package outbox

import (
	"context"

	pkgdb "starter-boilerplate/pkg/db"

	"github.com/uptrace/bun"
)

// Repository handles outbox table operations.
type Repository struct {
	db *bun.DB
}

func NewRepository(db *bun.DB) *Repository {
	return &Repository{db: db}
}

// Insert adds an entry to the outbox within the current tx (or fallback db).
func (r *Repository) Insert(ctx context.Context, entry *Entry) error {
	_, err := pkgdb.Conn(ctx, r.db).NewInsert().Model(entry).ExcludeColumn("id").Exec(ctx)
	return err
}

// FetchUnpublished returns up to limit unpublished entries, locking them for update.
func (r *Repository) FetchUnpublished(ctx context.Context, limit int) ([]Entry, error) {
	var entries []Entry
	err := pkgdb.Conn(ctx, r.db).NewSelect().
		Model(&entries).
		Where("published = FALSE").
		OrderExpr("id ASC").
		Limit(limit).
		For("UPDATE SKIP LOCKED").
		Scan(ctx)
	return entries, err
}

// MarkPublished sets published=TRUE for the given IDs.
func (r *Repository) MarkPublished(ctx context.Context, ids []int64) error {
	_, err := pkgdb.Conn(ctx, r.db).NewUpdate().
		Model((*Entry)(nil)).
		Set("published = TRUE").
		Where("id IN (?)", bun.In(ids)).
		Exec(ctx)
	return err
}
