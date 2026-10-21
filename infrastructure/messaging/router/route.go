package router

import (
	"context"
	"fmt"
	"go_project/infrastructure/messaging"
)

type router struct {
	routes map[string]messaging.Handler
}

func NewRouter() *router {
	return &router{
		routes: make(map[string]messaging.Handler),
	}
}

func (r *router) Handle(key string, fn messaging.Handler) {
	r.routes[key] = fn
}

// Dispatch выполняет поиск нужного хендлера и запуск
func (r *router) Dispatch(ctx context.Context, msg messaging.Message) error {
	// В качестве ключа маршрутизации выбираем msg.Topic
	// (или можно доставать event_type из msg.Headers / msg.Body)
	handler, exists := r.routes[msg.Topic]
	if !exists {
		return fmt.Errorf("no handler registered for route: %s", msg.Topic)
	}

	return handler(ctx, msg)
}
