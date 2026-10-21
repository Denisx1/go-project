package events

import (
	"context"
	messaging "go_project/infrastructure/messaging"
	"go_project/infrastructure/messaging/router"
	"go_project/order/internal/core/domain"
)

type OrderService interface {
	CreateOrder(
		ctx context.Context,
		order domain.Order,
	) error
}
type Router interface {
	Handle(key string, fn messaging.Handler)
	Dispatch(ctx context.Context, msg messaging.Message) error
}
type Handler struct {
	router       Router
	OrderService OrderService
}

func NewHandler(orderService OrderService) *Handler {
	h := &Handler{
		router:       router.NewRouter(),
		OrderService: orderService,
	}
	h.router.Handle("order.created", h.handleOrderCreated)
	return h
}

func (h *Handler) Handle(ctx context.Context, msg messaging.Message) error {
	return h.router.Dispatch(ctx, msg)
}
