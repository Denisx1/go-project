package postgres_pool

import (
	"context"

	"time"
)

// вот тут вот плозо чтоинерфейс зависит от библтотеки pgx
type Pool interface {
	Close()

	OpTimeout() time.Duration
	Exec(ctx context.Context, sql string, arguments ...any) (CommandTag, error)
	Query(ctx context.Context, sql string, args ...any) (Rows, error)
	QueryRow(ctx context.Context, sql string, args ...any) Row
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
