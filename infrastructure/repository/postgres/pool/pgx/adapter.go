package pgx_pool

import (
	"context"
	"errors"
	"fmt"
	error_type "go_project/infrastructure/error"
	postgres_pool "go_project/infrastructure/repository/postgres/pool"
	"io"
	"net"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

// errors.IS смотрит ощибки либо является либо сожержит в себе обернутую какую то конкретнцю ошибку которая назодится в переменно
// errors.As смотри либо получения ошибка либо в ее обернутых ощибказ есть структура конкретного переданого типа  и если есть то из за того что передаем по указателюб то она наполнится всеми полями что были заданы самой библтотекой pgx

type pgxCommandTag struct{ tag pgconn.CommandTag }

func (t pgxCommandTag) RowsAffected() int64 { return t.tag.RowsAffected() }

// pgxRow: cancel вызывается после Scan, иначе контекст с таймаутом утечёт.
type pgxRow struct {
	row    pgx.Row
	cancel context.CancelFunc
}

func (r pgxRow) Scan(dest ...any) error {
	defer r.cancel()
	return mapErrors(r.row.Scan(dest...))
}

// pgxRows: cancel вызывается в Close, ошибки маппятся.
type pgxRows struct {
	pgx.Rows
	cancel context.CancelFunc
}

func (r pgxRows) Scan(dest ...any) error { return mapErrors(r.Rows.Scan(dest...)) }
func (r pgxRows) Err() error             { return mapErrors(r.Rows.Err()) }
func (r pgxRows) Close() {
	r.Rows.Close()
	r.cancel()
}

type pgxTxAdapter struct{ tx pgx.Tx }

func (t *pgxTxAdapter) Exec(ctx context.Context, sql string, args ...any) (postgres_pool.CommandTag, error) {
	tag, err := t.tx.Exec(ctx, sql, args...)
	if err != nil {
		return nil, mapErrors(err)
	}
	return pgxCommandTag{tag: tag}, nil
}

func (t *pgxTxAdapter) Query(ctx context.Context, sql string, args ...any) (postgres_pool.Rows, error) {
	rows, err := t.tx.Query(ctx, sql, args...)
	if err != nil {
		return nil, mapErrors(err)
	}
	return pgxRows{Rows: rows, cancel: func() {}}, nil
}

func (t *pgxTxAdapter) QueryRow(ctx context.Context, sql string, args ...any) postgres_pool.Row {
	return pgxRow{row: t.tx.QueryRow(ctx, sql, args...), cancel: func() {}}
}

func (t *pgxTxAdapter) Commit(ctx context.Context) error   { return mapErrors(t.tx.Commit(ctx)) }
func (t *pgxTxAdapter) Rollback(ctx context.Context) error { return mapErrors(t.tx.Rollback(ctx)) }
func mapErrors(err error) error {
	if err == nil {
		return nil
	}
	const (
		// Полные коды
		codeViolatesCheck        = "23514"
		codeViolationNotNull     = "23502"
		codeUniqueViolation      = "23505"
		codeForeignKeyViolation  = "23503"
		codeSerializationFailure = "40001"
		codeDeadlockDetected     = "40P01"
		codeTimeoutExceeded      = "57014"
	)

	if errors.Is(err, pgx.ErrNoRows) {
		return postgres_pool.ErrNoRows
	}

	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		switch pgErr.Code {
		case codeUniqueViolation: // unique_violation
			return postgres_pool.ErrViolatesUnique
		case codeForeignKeyViolation: // foreign_key_violation
			return postgres_pool.ErrViolatesForeignKey
		case codeViolationNotNull: // not_null_violation
			return postgres_pool.ErrViolatesNotNull
		case codeViolatesCheck: // check_violation
			return postgres_pool.ErrViolatesCheck
		case codeSerializationFailure: // serialization_failure
			return postgres_pool.ErrSerializationFailure
		case codeDeadlockDetected: // deadlock_detected
			return postgres_pool.ErrDeadlockDetected
		case codeTimeoutExceeded: // query_canceled / statement_timeout
			return postgres_pool.ErrTimeoutExceeded
		}
	}
	// ошибки сети без PgError: connection refused, EOF, reset
	var netErr net.Error
	if errors.As(err, &netErr) || errors.Is(err, io.EOF) {
		return fmt.Errorf("%w: %w", error_type.ErrTransient, err)
	}

	return fmt.Errorf("%w: %w", postgres_pool.ErrUnknown, err)
}
