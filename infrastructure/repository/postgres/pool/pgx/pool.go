package pgx_pool

import (
	"context"
	"errors"
	"fmt"
	postgres_pool "go_project/infrastructure/repository/postgres/pool"
	"time"

	"github.com/bytedance/sonic"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Pool struct {
	*pgxpool.Pool
	opTimeout time.Duration
	txTimeout time.Duration
}

func NewPool(
	ctx context.Context,
	config Config,
) (*Pool, error) {
	connectionString := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		config.User,
		config.Password,
		config.Host,
		config.Port,
		config.Database,
	)
	pgxconfig, err := pgxpool.ParseConfig(connectionString)
	if err != nil {
		return nil, fmt.Errorf("parse pgxconfig: %w", err)
	}

	connectCtx, cancel := context.WithTimeout(ctx, config.ConnectTimeout)
	defer cancel()

	pgxconfig.AfterConnect = func(ctx context.Context, conn *pgx.Conn) error {
		codec := &pgtype.JSONCodec{
			Marshal:   sonic.Marshal,
			Unmarshal: sonic.Unmarshal,
		}
		conn.TypeMap().RegisterType(&pgtype.Type{Name: "json", OID: pgtype.JSONOID, Codec: codec})
		conn.TypeMap().RegisterType(&pgtype.Type{Name: "jsonb", OID: pgtype.JSONBOID, Codec: codec})
		return nil
	}
	pool, err := pgxpool.NewWithConfig(connectCtx, pgxconfig)
	if err != nil {
		return nil, fmt.Errorf("create pgxpool: %w", err)
	}

	if err := pool.Ping(connectCtx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("pgxpool ping: %w", err)
	}
	return &Pool{
		Pool:      pool,
		opTimeout: config.OpTimeout,
		txTimeout: config.TxTimeout,
	}, nil
}

func (p *Pool) OpTimeout() time.Duration { return p.opTimeout }

func (p *Pool) Close() { p.Pool.Close() }

func (p *Pool) withTimeout(ctx context.Context) (context.Context, context.CancelFunc) {
	if p.opTimeout <= 0 {
		return ctx, func() {}
	}
	if deadline, ok := ctx.Deadline(); ok && time.Until(deadline) < p.opTimeout {
		return ctx, func() {}
	}
	return context.WithTimeout(ctx, p.opTimeout)
}

func txFromContext(ctx context.Context) (postgres_pool.Tx, bool) {
	tx, ok := ctx.Value(postgres_pool.TxKey{}).(postgres_pool.Tx)
	return tx, ok
}

func (p *Pool) Exec(
	ctx context.Context,
	sql string,
	args ...any,
) (postgres_pool.CommandTag, error) {
	if tx, ok := txFromContext(ctx); ok {
		return tx.Exec(ctx, sql, args...)
	}
	txCtx, cancel := p.withDeadline(ctx, p.txTimeout)
	defer cancel()
	tag, err := p.Pool.Exec(txCtx, sql, args...)
	if err != nil {
		return nil, err
	}
	return pgxCommandTag{tag: tag}, nil
}
func (p *Pool) Query(
	ctx context.Context,
	sql string,
	args ...any) (postgres_pool.Rows, error) {
	if tx, ok := txFromContext(ctx); ok {
		return tx.Query(ctx, sql, args...)
	}
	txCtx, cancel := p.withDeadline(ctx, p.txTimeout)

	rows, err := p.Pool.Query(txCtx, sql, args...)
	if err != nil {
		cancel()
		return nil, err
	}
	return pgxRows{Rows: rows, cancel: cancel}, nil
}
func (p *Pool) QueryRow(
	ctx context.Context,
	sql string,
	args ...any) postgres_pool.Row {
	if tx, ok := txFromContext(ctx); ok {
		return tx.QueryRow(ctx, sql, args...)
	}
	txCtx, cancel := p.withDeadline(ctx, p.txTimeout)
	return pgxRow{row: p.Pool.QueryRow(txCtx, sql, args...), cancel: cancel}
}

func (p *Pool) withTx(
	ctx context.Context,
	opts pgx.TxOptions,
	fn func(ctx context.Context) error,
) (err error) {
	// Вложенный вызов: уже в транзакции, работаем в ней
	if _, ok := txFromContext(ctx); ok {
		return fn(ctx)
	}

	// Таймаут на всю транзакцию
	txCtx, cancel := p.withDeadline(ctx, p.txTimeout)
	defer cancel()

	pgxTx, err := p.Pool.BeginTx(txCtx, opts)
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}
	tx := &pgxTxAdapter{tx: pgxTx}
	ctxWithTx := context.WithValue(txCtx, postgres_pool.TxKey{}, postgres_pool.Tx(tx))

	defer func() {
		if r := recover(); r != nil {
			rollback(ctx, pgxTx)
			panic(r)
		}
		if err != nil {
			rollback(ctx, pgxTx)
			return
		}
		if cErr := pgxTx.Commit(txCtx); cErr != nil {
			err = fmt.Errorf("commit tx: %w", cErr)
		}
	}()

	return fn(ctxWithTx)
}
func (p *Pool) withDeadline(
	ctx context.Context,
	d time.Duration,
) (context.Context, context.CancelFunc) {
	if d <= 0 {
		return ctx, func() {}
	}
	if deadline, ok := ctx.Deadline(); ok && time.Until(deadline) < d {
		return ctx, func() {}
	}
	return context.WithTimeout(ctx, d)
}

func (p *Pool) WithTx(ctx context.Context, fn func(ctx context.Context) error) error {
	return p.withTx(ctx, pgx.TxOptions{}, fn)
}
func (p *Pool) WithTxSerializable(ctx context.Context, fn func(ctx context.Context) error) error {
	const maxRetries = 3
	opts := pgx.TxOptions{IsoLevel: pgx.Serializable}

	var err error
	for i := 0; i <= maxRetries; i++ {
		if err = p.withTx(ctx, opts, fn); err == nil || !isSerializationFailure(err) {
			return err
		}
	}
	return fmt.Errorf("tx failed after %d retries: %w", maxRetries, err)
}

// Rollback на «свежем» контексте: если txCtx уже истёк по таймауту,
// откат с ним не выполнится и соединение зависнет в незакрытой транзакции.
func rollback(ctx context.Context, tx pgx.Tx) {
	rbCtx, cancel := context.WithTimeout(context.WithoutCancel(ctx), 5*time.Second)
	defer cancel()
	_ = tx.Rollback(rbCtx)
}

func isSerializationFailure(err error) bool {
	var pgErr *pgconn.PgError
	return errors.As(err, &pgErr) && pgErr.Code == "40001"
}
