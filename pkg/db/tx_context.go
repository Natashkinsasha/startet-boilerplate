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

// TxFromCtx returns the bun.Tx stored in ctx, if any.
func TxFromCtx(ctx context.Context) (bun.Tx, bool) {
	tx, ok := ctx.Value(txKey{}).(bun.Tx)
	return tx, ok
}

// Conn returns the active connection from ctx: transaction first, then the fallback *bun.DB.
func Conn(ctx context.Context, fallback *bun.DB) bun.IDB {
	if tx, ok := TxFromCtx(ctx); ok {
		return tx
	}
	return fallback
}
