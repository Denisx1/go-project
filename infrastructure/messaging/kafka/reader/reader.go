package kafka_reader

import (
	"context"
	"fmt"
	messaging "go_project/infrastructure/messaging"
	"go_project/infrastructure/messaging/component/consumer"

	kafkago "github.com/segmentio/kafka-go"
)

type Reader struct {
	r *kafkago.Reader
}

var _ consumer.Reader = (*Reader)(nil)

func NewReader(
	cfg ReaderConfig,
) (*Reader, error) {
	start := kafkago.FirstOffset
	if cfg.StartFromLatest {
		start = kafkago.LastOffset
	}

	return &Reader{
		r: kafkago.NewReader(kafkago.ReaderConfig{
			Brokers:     cfg.CommonConfig.Brokers,
			Dialer:      cfg.CommonConfig.Dialer(),
			GroupTopics: cfg.Topics,
			GroupID:     cfg.GroupID,
			StartOffset: start,
			MinBytes:    cfg.MinBytes,
			MaxBytes:    cfg.MaxBytes,
			MaxWait:     cfg.MaxWait,
			// CommitInterval не задаём (0): CommitMessages синхронный, коммитим вручную после обработки.
		})}, nil
}

func (s *Reader) FetchMessage(ctx context.Context) (messaging.Delivery, error) {
	km, err := s.r.FetchMessage(ctx)
	if err != nil {
		// ctx.Err() != nil или io.EOF (Reader закрыт) - Consumer решит, что это штатная остановка
		return messaging.Delivery{}, fmt.Errorf("kafka fetch: %w", err)
	}

	return messaging.Delivery{
		Message: s.toMessage(km),
		Commit: func(ctx context.Context) error {
			return s.commitMessage(ctx, km)
		},
	}, nil
}
func (s *Reader) commitMessage(
	ctx context.Context,
	km kafkago.Message,
) error {

	if err := s.r.CommitMessages(ctx, km); err != nil {
		if ctx.Err() != nil {
			return nil
		}
		return fmt.Errorf("kafka commit %s[%d]@%d: %w", km.Topic, km.Partition, km.Offset, err)
	}
	return nil
}

func (s *Reader) Close() error { return s.r.Close() }

func (s *Reader) toMessage(km kafkago.Message) messaging.Message {
	hs := make(map[string]string, len(km.Headers))
	for _, h := range km.Headers {
		hs[h.Key] = string(h.Value)
	}
	return messaging.Message{
		Topic:   km.Topic,
		Key:     km.Key,
		Body:    km.Value,
		Headers: hs,
	}
}

// func (s *Reader) Subscribe(
// 	ctx context.Context,
// 	h messaging.Handler,
// ) error {
// 	for {
// 		km, err := s.r.FetchMessage(ctx)
// 		if err != nil {
// 			if ctx.Err() != nil || errors.Is(err, io.EOF) { // io.EOF - Reader закрыт через Close
// 				return nil
// 			}
// 			return fmt.Errorf("kafka fetch (%s): %w", km.Topic, err)
// 		}

// 		if err := h(ctx, toMessage(km)); err != nil {
// 			if ctx.Err() != nil {
// 				return nil // не закоммичено -> будет прочитано заново
// 			}
// 			// Handler окончательно не справился (даже DLQ недоступен).
// 			// Не коммитим и останавливаемся: сообщение не теряется.
// 			return fmt.Errorf("kafka handle %s[%d]@%d: %w", km.Topic, km.Partition, km.Offset, err)
// 		}

// 		if err := s.r.CommitMessages(ctx, km); err != nil {
// 			if ctx.Err() != nil {
// 				return nil
// 			}
// 			return fmt.Errorf("kafka commit %s[%d]@%d: %w", km.Topic, km.Partition, km.Offset, err)
// 		}
// 	}
// }
