package postgres_pool

import (
	"errors"
)

var (
	ErrNoRows             = errors.New("no rows found")
	ErrViolatesUnique     = errors.New("violates unique constraint (duplicate key)")
	ErrViolatesForeignKey = errors.New("violates foreign key constraint")
	ErrViolatesNotNull    = errors.New("violates not null constraint")
	ErrViolatesCheck      = errors.New("violates check constraint")

	// --- 2. Конкурентный доступ и транзакции (Concurrency & TX) ---
	ErrSerializationFailure = errors.New("transaction serialization failure (optimistic lock conflict)")
	ErrDeadlockDetected     = errors.New("database deadlock detected")
	ErrTxClosed             = errors.New("transaction is already closed, committed, or rolled back")
	ErrTxBeginFailed        = errors.New("failed to begin transaction")
	ErrTxCommitFailed       = errors.New("failed to commit transaction")
	ErrTxRollbackFailed     = errors.New("failed to rollback transaction")

	// --- 3. Соединения, пул и тайм-ауты (Connection & Pool) ---
	ErrConnectionClosed = errors.New("database connection is closed")
	ErrPoolClosed       = errors.New("connection pool is closed")
	ErrTimeoutExceeded  = errors.New("database query timeout exceeded")

	ErrInvalidInput = errors.New("invalid input in db")

	ErrSyntax = errors.New("invalid sql")

	// --- 4. Общая / Неизвестная ошибка ---
	ErrUnknown = errors.New("unknown database error")
)
