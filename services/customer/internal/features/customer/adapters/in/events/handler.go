package events

import (
	"context"
	orderevents "go_project/contracts/events"
	messaging "go_project/infrastructure/messaging"
	"go_project/infrastructure/messaging/router"
)

type Router interface {
	Handle(key string, fn messaging.Handler)
	Dispatch(ctx context.Context, msg messaging.Message) error
}
type Handler struct {
	router Router
}

func NewHandler() *Handler {
	h := &Handler{
		router: router.NewRouter(),
	}
	h.router.Handle(orderevents.EventTypeCreated, h.createCustomer)
	return h
}

func (h *Handler) Handle(ctx context.Context, msg messaging.Message) error {
	return h.router.Dispatch(ctx, msg)
}
