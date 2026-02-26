package db

import (
	"context"
	"database/sql"

	"github.com/uptrace/bun"
)

// UoW abstracts transactional execution.
type UoW interface {
	Do(ctx context.Context, fn func(ctx context.Context) error, opts ...*sql.TxOptions) error
}

// UnitOfWork wraps *bun.DB and provides transactional execution.
type UnitOfWork struct {
	db *bun.DB
}

func NewUnitOfWork(db *bun.DB) *UnitOfWork {
	return &UnitOfWork{db: db}
}

// Do executes fn inside a database transaction. The transaction is stored
// in the child context via WithTx so Conn picks it up automatically.
func (u *UnitOfWork) Do(ctx context.Context, fn func(ctx context.Context) error, opts ...*sql.TxOptions) error {
	var txOpts *sql.TxOptions
	if len(opts) > 0 {
		txOpts = opts[0]
	}
	return u.db.RunInTx(ctx, txOpts, func(ctx context.Context, tx bun.Tx) error {
		return fn(WithTx(ctx, tx))
	})
}
