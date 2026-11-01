package events

import (
	"context"
	"fmt"
	logger "go_project/infrastructure/logging"
	messaging "go_project/infrastructure/messaging"
	"go_project/infrastructure/messaging/request"
	"go_project/order/internal/core/domain"
)

type CreateOrderEvent struct {
	Title string `json:"title" validate:"required,min=10,max=100"`
	Text  string `json:"text" validate:"required,min=10,max=100"`
}

func (h *Handler) handleOrderCreated(
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
	var event CreateOrderEvent
	if err := request.DecodeAndValidate(msg, &event); err != nil {
		return fmt.Errorf("failed to decode and validate request: %w", err)
	}
	orderDomain := domainFromDTO(event)
	_, err := h.OrderService.CreateOrder(ctx, orderDomain)
	if err != nil {
		return fmt.Errorf("failed to handle order: %w", err)
	}
	return nil

}

func domainFromDTO(dto CreateOrderEvent) domain.Order {
	return domain.NewOrderUninitialized(dto.Title, dto.Text)
}
