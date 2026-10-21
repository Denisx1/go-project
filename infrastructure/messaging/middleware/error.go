package middleware

import (
	"errors"
	"fmt"
)

// nonRetryableError — ошибка, которую нет смысла повторять (битый JSON, валидация и т.д.)
type nonRetryableError struct{ err error }

func (e *nonRetryableError) Error() string { return e.err.Error() }
func (e *nonRetryableError) Unwrap() error { return e.err }

// NonRetryable оборачивает ошибку: ретрай пропускается, сообщение сразу уходит в DLQ.
func NonRetryable(err error) error {
	if err == nil {
		return nil
	}
	return &nonRetryableError{err: err}
}

func IsNonRetryable(err error) bool {
	var nr *nonRetryableError
	return errors.As(err, &nr)
}

// exhaustedError — ретраи закончились; хранит число попыток для DLQ.
type exhaustedError struct {
	err      error
	attempts int
}

func (e *exhaustedError) Error() string {
	return fmt.Sprintf("retries exhausted after %d attempts: %v", e.attempts, e.err)
}
func (e *exhaustedError) Unwrap() error { return e.err }
