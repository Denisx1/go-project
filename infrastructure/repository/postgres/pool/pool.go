package postgres_pool

import (
	"context"
	"time"
)

// вот тут вот плозо чтоинерфейс зависит от библтотеки pgx

type Pool interface {
	Exec(ctx context.Context, sql string, args ...any) (CommandTag, error)
	Query(ctx context.Context, sql string, args ...any) (Rows, error)
	QueryRow(ctx context.Context, sql string, args ...any) Row
	OpTimeout() time.Duration
}

type Tx interface {
	Exec(ctx context.Context, sql string, args ...any) (CommandTag, error)
	Query(ctx context.Context, sql string, args ...any) (Rows, error)
	QueryRow(ctx context.Context, sql string, args ...any) Row
	Commit(ctx context.Context) error
	Rollback(ctx context.Context) error
}
type CommandTag interface {
	RowsAffected() int64
}
type Row interface {
	Scan(dest ...any) error
}
type Rows interface {
	Close()
	Err() error
	Next() bool
	Scan(dest ...any) error
}
type TxKey struct{}
