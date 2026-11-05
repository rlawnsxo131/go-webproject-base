package repository

import (
	"database/sql"
	"log/slog"
)

// type QueryOption interface {
// 	apply(*queryOptions)
// }

type queryOptions struct {
	WithLock bool
	Tx       *sql.Tx
}

// verify-interface-compliance
// @link https://github.com/uber-go/guide/blob/master/style.md#verify-interface-compliance
// var _ QueryOption = (*queryOptionFunc)(nil)

type queryOptionFunc func(*queryOptions)

func (f queryOptionFunc) apply(o *queryOptions) {
	f(o)
}

func WithLock() queryOptionFunc {
	return queryOptionFunc(func(o *queryOptions) {
		o.WithLock = true
	})
}

func WithTx(tx *sql.Tx) queryOptionFunc {
	return queryOptionFunc(func(o *queryOptions) {
		o.Tx = tx
	})
}

func parseOptions(opts []queryOptionFunc) *queryOptions {
	qOpts := &queryOptions{WithLock: false, Tx: nil}
	for _, opt := range opts {
		opt.apply(qOpts)
	}
	return qOpts
}

func debug(query string, args any) {
	slog.Debug("[SQL]",
		slog.String("query", query),
		slog.Any("args", args),
	)
}
