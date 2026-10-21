package request

import (
	"encoding/json"
	"fmt"
	error_type "go_project/infrastructure/error"
	"go_project/infrastructure/messaging"

	"gopkg.in/go-playground/validator.v9"
)

var requestValidater = validator.New()

type Validatable interface {
	Validate() error
}

func DecodeAndValidate[T any](msg messaging.Message, dst *T) error {
	if err := json.Unmarshal(msg.Body, dst); err != nil {
		return fmt.Errorf("%w: decode: %w", error_type.ErrPermanent, err)
	}
	var err error
	if v, ok := any(dst).(Validatable); ok {
		err = v.Validate() // кастомные правила
	} else {
		err = requestValidater.Struct(dst)
	}
	if err != nil {
		return fmt.Errorf("%w: validate: %w", error_type.ErrPermanent, err)
	}

	return nil
}
