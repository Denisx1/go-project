package error_type

import "errors"

var (
	ErrTransient = errors.New("transient") // можно повторить
	ErrPermanent = errors.New("permanent") // повтор бесполезен
)

func IsRetryable(err error) bool {
	return !errors.Is(err, ErrPermanent)
}

var (
	ErrNotFound        = errors.New("not found")
	ErrInvalidArgument = errors.New("invalid argument")
	ErrConflict        = errors.New("conflict")
	ErrAlreadyExist    = errors.New("already exists")
)
