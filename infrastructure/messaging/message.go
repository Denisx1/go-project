package messaging

import (
	"context"
)

type Message struct {
	Topic   string
	Key     []byte
	Body    []byte
	Headers map[string]string
}
type Delivery struct {
	Message Message
	Commit  func(ctx context.Context) error
}

type Handler func(ctx context.Context, msg Message) error

type Consumer interface {
	Start(ctx context.Context) error
	Stop(ctx context.Context) error
}

type Producer interface {
	Publish(ctx context.Context, msg Message) error
	PublishDLT(
		ctx context.Context,
		msg Message,
		cause error,
		attempts int,
	) error
}
