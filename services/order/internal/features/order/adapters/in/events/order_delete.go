package events

import (
	"context"
	"fmt"
	logger "go_project/infrastructure/logging"
	messaging "go_project/infrastructure/messaging"
	"go_project/infrastructure/messaging/request"
)

type DeleteOrderEvent struct {
	OrderId string `json:"order_id"`
}

func (h *Handler) handleDeleteOrder(
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
	var event DeleteOrderEvent
	if err := request.DecodeAndValidate(msg, &event); err != nil {
		return fmt.Errorf("failed to decode and validate request: %w", err)
	}
	if err := h.OrderService.DeleteOrder(ctx, event.OrderId); err != nil {
		return fmt.Errorf("failed to handle order: %w", err)
	}
	return nil

}
