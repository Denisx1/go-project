package postgres_pool

import (
	"context"
	"time"
)

// вот тут вот плозо чтоинерфейс зависит от библтотеки pgx

// Executor: общее, что умеют и пул, и транзакция.
type Executor interface {
	Exec(ctx context.Context, sql string, args ...any) (CommandTag, error)
	Query(ctx context.Context, sql string, args ...any) (Rows, error)
	QueryRow(ctx context.Context, sql string, args ...any) Row
}

// Pool: запросы вне транзакции.
type Pool interface {
	Executor
	OpTimeout() time.Duration
}

// Tx: запросы внутри транзакции.
type Tx interface {
	Executor
	Commit(ctx context.Context) error
	Rollback(ctx context.Context) error
}

// TxManager: запуск транзакций.
type TxManager interface {
	WithTx(ctx context.Context, fn func(ctx context.Context) error) error
	WithTxSerializable(ctx context.Context, fn func(ctx context.Context) error) error
	WithTxRepeatableRead(ctx context.Context, fn func(ctx context.Context) error) error
}

// Lifecycle: служебные методы.
type Lifecycle interface {
	Ping(ctx context.Context) error
	Close()
}

// Database: всё вместе, это возвращает конструктор.
type Database interface {
	Pool
	TxManager
	Lifecycle
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
