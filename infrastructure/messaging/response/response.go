package response

import (
	"errors"
	"fmt"
	error_type "go_project/infrastructure/error"
	logger "go_project/infrastructure/logging"
	"go_project/infrastructure/messaging"
)

type ResultHandler struct {
	log logger.Logger
	msg messaging.Message
}

func NewResultHandler(
	log logger.Logger,
	msg messaging.Message,
) *ResultHandler {
	return &ResultHandler{log: log, msg: msg}
}

func (h *ResultHandler) ErrorResponse(err error, msg string) error {
	if err == nil {
		return nil
	}

	var (
		marker  error
		logFunc func(string, ...any)
	)

	switch {
	case errors.Is(err, error_type.ErrConflict):
		h.log.Info(msg+": duplicate, treated as success", append(h.fields(), "reason", err)...)
		return nil

	case errors.Is(err, error_type.ErrNotFound),
		errors.Is(err, error_type.ErrInvalidArgument):
		marker, logFunc = error_type.ErrPermanent, h.log.Warn

	default:
		marker, logFunc = error_type.ErrTransient, h.log.Error
	}

	logFunc(msg, append(h.fields(), "error", err)...)

	// сообщение + маркер + исходная причина, цепочка errors.Is сохраняется
	if errors.Is(err, marker) {
		return fmt.Errorf("%s: %w", msg, err) // маркер уже в цепочке, второй раз не добавляем
	}
	return fmt.Errorf("%s: %w: %w", msg, marker, err)

}
func (h *ResultHandler) fields() []any {
	return []any{
		"topic", h.msg.Topic,
		"key", string(h.msg.Key),
	}
}
