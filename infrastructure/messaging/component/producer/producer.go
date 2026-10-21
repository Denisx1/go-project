package producer

import (
	"context"
	"go_project/infrastructure/messaging"
	"strconv"
)

type Writer interface {
	Publish(ctx context.Context, msg messaging.Message) error
	GetDltSuffix() string
	Close() error
}

type Producer struct {
	pub Writer
}

func NewProducer(pub Writer) *Producer {
	return &Producer{pub: pub}
}

func (p *Producer) Publish(ctx context.Context, msg messaging.Message) error {

	return p.pub.Publish(ctx, msg)
}

func (p *Producer) PublishDLQ(
	ctx context.Context,
	msg messaging.Message,
	cause error,
	attempts int,
) error {
	headers := make(map[string]string, len(msg.Headers)+3)
	suffix := p.pub.GetDltSuffix()

	dltTopic := msg.Topic + suffix
	for k, v := range msg.Headers { // копируем, чтобы не мутировать исходную map
		headers[k] = v
	}
	headers["x-original-topic"] = msg.Topic
	headers["x-error"] = cause.Error()
	headers["x-attempts"] = strconv.Itoa(attempts)

	msg.Topic = dltTopic
	msg.Headers = headers
	return p.Publish(ctx, msg)
}

func (p *Producer) Close() error { return p.pub.Close() }
