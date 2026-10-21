package kafka_writter

import (
	"context"

	"go_project/infrastructure/messaging"

	kafkago "github.com/segmentio/kafka-go"
)

type Writer struct {
	w   *kafkago.Writer
	cfg WriterConfig
}

func NewWritter(cfg WriterConfig) (*Writer, error) {

	writer := &kafkago.Writer{
		Addr:     kafkago.TCP(cfg.CommonConfig.Brokers...),
		Balancer: &kafkago.Hash{},
		// BatchTimeout:           cfg.BatchTimeout,
		WriteTimeout:           cfg.WriteTimeout,
		AllowAutoTopicCreation: cfg.AllowAutoTopicCreation,
		RequiredAcks:           kafkago.RequireAll,
		BatchSize:              1, // отправлять сразу при поступлении 1 сообщения
		Async:                  false,
		MaxAttempts:            3,
	}

	return &Writer{w: writer, cfg: cfg}, nil
}
func (w *Writer) GetDltSuffix() string {
	if w.cfg.DLTSuffix == "" {
		return ".DLT"
	}
	return w.cfg.DLTSuffix
}
func (p *Writer) Publish(
	ctx context.Context,
	msg messaging.Message,
) error {
	newHeaders := make([]kafkago.Header, 0, len(msg.Headers))
	for key, value := range msg.Headers {
		newHeaders = append(newHeaders, kafkago.Header{
			Key:   key,
			Value: []byte(value),
		})
	}

	return p.w.WriteMessages(ctx, kafkago.Message{
		Topic:   msg.Topic,
		Key:     msg.Key,
		Value:   msg.Body,
		Headers: newHeaders,
	})
}
func (p *Writer) Close() error {
	return p.w.Close()
}
