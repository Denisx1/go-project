package postgres_pool

import (
	"errors"
	"fmt"
	error_type "go_project/infrastructure/error"
)

var (
	ErrViolatesForeignKey = fmt.Errorf("%w: violates foreign key", error_type.ErrPermanent)
	ErrUniqueViolation    = fmt.Errorf("%w: unique violation", error_type.ErrPermanent)

	// повтор может помочь
	ErrConnection = fmt.Errorf("%w: connection error", error_type.ErrTransient)
	ErrDeadlock   = fmt.Errorf("%w: deadlock", error_type.ErrTransient)
	ErrTimeout    = fmt.Errorf("%w: timeout", error_type.ErrTransient)

	// без маркера
	ErrNoRows  = errors.New("no rows")
	ErrUnknown = errors.New("unknown")
)
