package events

import (
	"context"
	logger "go_project/infrastructure/logging"
	messaging "go_project/infrastructure/messaging"
)

func (h *Handler) createCustomer(
	ctx context.Context,
	msg messaging.Message,
) error {
	log := logger.FromContext(ctx)
	// resposeHandler := response.NewResultHandler(log, msg)
	log.Info("handling message",
		"topic", msg.Topic,
		"key", string(msg.Key),
		"body_size", len(msg.Body),
		"headers", msg.Headers,
	)
	return nil
}
