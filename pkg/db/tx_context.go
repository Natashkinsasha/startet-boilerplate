package db

import (
	"context"

	"github.com/uptrace/bun"
)

type txKey struct{}

// WithTx stores a bun.Tx in the context.
func WithTx(ctx context.Context, tx bun.Tx) context.Context {
	return context.WithValue(ctx, txKey{}, tx)
}

// Conn returns the bun.Tx stored in ctx, or falls back to the given *bun.DB.
func Conn(ctx context.Context, fallback *bun.DB) bun.IDB {
	if tx, ok := ctx.Value(txKey{}).(bun.Tx); ok {
		return tx
	}
	return fallback
}

// RunInTx executes fn inside a database transaction.
// The transaction is stored in the context via WithTx so repositories
// can pick it up automatically through Conn.
func RunInTx(ctx context.Context, db *bun.DB, fn func(ctx context.Context) error) error {
	return db.RunInTx(ctx, nil, func(ctx context.Context, tx bun.Tx) error {
		return fn(WithTx(ctx, tx))
	})
}
