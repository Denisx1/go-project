package grpcrequest

import (
	"fmt"
	"strings"
	"unicode/utf8"

	error_type "go_project/infrastructure/error"

	"github.com/google/uuid"
)

type Validator struct {
	errs []string
}

func New() *Validator { return &Validator{} }

// Add добавляет свою ошибку (для нестандартных проверок).
func (v *Validator) Add(field, msg string) *Validator {
	v.errs = append(v.errs, field+": "+msg)
	return v
}

func (v *Validator) Required(field, value string) *Validator {
	if strings.TrimSpace(value) == "" {
		v.Add(field, "is required")
	}
	return v
}

func (v *Validator) MaxLen(field, value string, max int) *Validator {
	if utf8.RuneCountInString(value) > max {
		v.Add(field, fmt.Sprintf("must be at most %d characters", max))
	}
	return v
}

func (v *Validator) UUID(field, value string) *Validator {
	if _, err := uuid.Parse(value); err != nil {
		v.Add(field, "must be a valid UUID")
	}
	return v
}

func (v *Validator) Positive(field string, value int64) *Validator {
	if value <= 0 {
		v.Add(field, "must be greater than 0")
	}
	return v
}

// OptionalNotEmpty: nil = поле не передано, проверка пропускается.
func (v *Validator) OptionalNotEmpty(field string, value *string, max int) *Validator {
	if value == nil {
		return v
	}
	return v.Required(field, *value).MaxLen(field, *value, max)
}

// Err возвращает nil или ошибку, обёрнутую в ErrInvalidArgument.
func (v *Validator) Err() error {
	if len(v.errs) == 0 {
		return nil
	}
	return fmt.Errorf("%w: %s", error_type.ErrInvalidArgument, strings.Join(v.errs, "; "))
}
