package pgx_pool

import (
	"context"
	"errors"
	"fmt"
	error_type "go_project/infrastructure/error"
	postgres_pool "go_project/infrastructure/repository/postgres/pool"
	"io"
	"net"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

// errors.IS смотрит ощибки либо является либо сожержит в себе обернутую какую то конкретнцю ошибку которая назодится в переменно
// errors.As смотри либо получения ошибка либо в ее обернутых ощибказ есть структура конкретного переданого типа  и если есть то из за того что передаем по указателюб то она наполнится всеми полями что были заданы самой библтотекой pgx

const (
	// Полные коды
	codeUniqueViolation      = "23505"
	codeForeignKeyViolation  = "23503"
	codeSerializationFailure = "40001"
	codeDeadlockDetected     = "40P01"
	codeAdminShutdown        = "57P01"
	codeTooManyConnections   = "53300"

	// Префиксы классов
	classConnectionException = "08"
	classIntegrityViolation  = "23"
	classDataException       = "22"
	classSyntaxOrAccess      = "42"
)

type pgxRow struct{ pgx.Row }

type pgxRows struct{ pgx.Rows }

func (r pgxRow) Scan(dest ...any) error {

	err := r.Row.Scan(dest...)
	if err != nil {
		return mapErrors(err)
	}
	return nil
}

type pgxCommandTag struct{ pgconn.CommandTag }

func isTransientCode(code string) bool {
	switch code {
	case codeSerializationFailure,
		codeDeadlockDetected,
		codeAdminShutdown,
		codeTooManyConnections:
		return true
	}
	return strings.HasPrefix(code, classConnectionException)
}

func mapErrors(err error) error {
	if errors.Is(err, pgx.ErrNoRows) {
		return postgres_pool.ErrNoRows
	}

	// контекст: таймаут или отмена
	if errors.Is(err, context.DeadlineExceeded) {
		return fmt.Errorf("%w: %w", error_type.ErrTransient, err)
	}

	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		switch {
		case pgErr.Code == codeUniqueViolation:
			return fmt.Errorf("%w: %w", postgres_pool.ErrUniqueViolation, err)
		case strings.HasPrefix(pgErr.Code, classSyntaxOrAccess):
			// баг в SQL: повтор не поможет
			return fmt.Errorf("%w: %w", error_type.ErrPermanent, err)
		case pgErr.Code == codeForeignKeyViolation:
			return fmt.Errorf("%w: %w", postgres_pool.ErrViolatesForeignKey, err)
		case isTransientCode(pgErr.Code):
			return fmt.Errorf("%w: %w", error_type.ErrTransient, err)
		case strings.HasPrefix(pgErr.Code, classIntegrityViolation),
			strings.HasPrefix(pgErr.Code, classDataException):
			// остальные нарушения ограничений и ошибки данных: повтор не поможет
			return fmt.Errorf("%w: %w", error_type.ErrPermanent, err)
		}
	}

	// ошибки сети без PgError: connection refused, EOF, reset
	var netErr net.Error
	if errors.As(err, &netErr) || errors.Is(err, io.EOF) {
		return fmt.Errorf("%w: %w", error_type.ErrTransient, err)
	}

	return fmt.Errorf("%w: %w", postgres_pool.ErrUnknown, err)
}
