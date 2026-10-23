package middleware

import (
	"context"
	"errors"
	"fmt"
	error_type "go_project/infrastructure/error"
	logger "go_project/infrastructure/logging"
	messaging "go_project/infrastructure/messaging"
	"go_project/infrastructure/messaging/component/producer"
	"runtime/debug"
	"time"

	"github.com/google/uuid"
)

const HeaderRequestID = "x-request-id"

type requestIDKey struct{}

// ---------- RequestID ----------

func withRequestID() Middleware {
	return func(next messaging.Handler) messaging.Handler {
		return func(ctx context.Context, msg messaging.Message) error {
			id := msg.Headers[HeaderRequestID]
			if id == "" {
				id = uuid.NewString()
			}
			ctx = context.WithValue(ctx, requestIDKey{}, id)
			return next(ctx, msg)
		}
	}
}

// ---------- Log ----------

func withLog(log logger.Logger) Middleware {
	return func(next messaging.Handler) messaging.Handler {
		return func(ctx context.Context, msg messaging.Message) error {
			// 1. Достаем request_id, который положила предыдущая middleware (withRequestID)
			requestId, _ := ctx.Value(requestIDKey{}).(string)

			// 2. Создаем расширенный логгер
			extendedLog := log.With(
				"request_id", requestId,
				"size", len(msg.Body),
			)

			newCtx := logger.ToContext(ctx, extendedLog)

			return next(newCtx, msg)
		}
	}
}

// ---------- Trace ----------

func withTrace() Middleware {
	return func(next messaging.Handler) messaging.Handler {
		return func(ctx context.Context, msg messaging.Message) error {
			log := logger.FromContext(ctx)
			before := time.Now()

			log.Debug(
				">>> Incoming Async request",
				"Topic", msg.Topic,
				"StartTime", before.UTC(),
			)
			next(ctx, msg)
			log.Debug(
				"<<< Done async request",
				"Latency", time.Since(before).String(),
			)
			return nil
		}
	}
}

// ---------- Retry ----------

func withRetry(attempts int) Middleware {
	if attempts < 1 {
		attempts = 1
	}
	const baseBackoff = 200 * time.Millisecond

	return func(next messaging.Handler) messaging.Handler {
		return func(ctx context.Context, msg messaging.Message) error {
			log := logger.FromContext(ctx)

			var err error
			for i := 0; i < attempts; i++ {
				if err = next(ctx, msg); err == nil {
					if i > 0 {
						log.Info("message succeeded after retry", "attempt", i+1)
					}
					return nil
				}

				// 1. Решение читаем из маркеров, а не из IsNonRetryable
				if !error_type.IsRetryable(err) {
					log.Warn("non-retryable error, skipping retries",
						"attempt", i+1, "error", err)
					return err // без exhaustedError: попытки не исчерпаны
				}

				// 2. Контекст отменён (shutdown): не ретраим и не отправляем в DLQ
				if ctx.Err() != nil {
					return ctx.Err()
				}

				// 3. Последняя попытка
				if i == attempts-1 {
					log.Error("retry attempts exhausted",
						"attempt", i+1, "max_attempts", attempts, "error", err)
					break
				}

				delay := baseBackoff * time.Duration(1<<i) // 200ms, 400ms, 800ms
				log.Warn("retry attempt failed, will retry",
					"attempt", i+1,
					"max_attempts", attempts,
					"next_retry_in", delay.String(),
					"error", err,
				)

				t := time.NewTimer(delay)
				select {
				case <-ctx.Done():
					t.Stop()
					return ctx.Err()
				case <-t.C:
				}
			}
			return &exhaustedError{err: err, attempts: attempts}
		}
	}
}
func withDLQ(writer producer.Writer) Middleware {
	dltProducer := producer.NewProducer(writer) // создаём один раз, а не на каждое сообщение

	return func(next messaging.Handler) messaging.Handler {
		return func(ctx context.Context, msg messaging.Message) error {
			err := next(ctx, msg)
			if err == nil {
				return nil
			}

			// Контекст отменён (graceful shutdown) — это не вина сообщения,
			// не отправляем в DLQ, пусть переобработается после рестарта.
			if ctx.Err() != nil {
				return err
			}

			log := logger.FromContext(ctx)

			attempts := 1
			var ex *exhaustedError
			if errors.As(err, &ex) {
				attempts = ex.attempts
			}

			// DLQ-публикация не должна зависеть от отмены исходного контекста
			dlqCtx, cancel := context.WithTimeout(context.WithoutCancel(ctx), 15*time.Second)
			defer cancel()

			if pubErr := dltProducer.PublishDLQ(dlqCtx, msg, err, attempts); pubErr != nil {
				log.Error("failed to publish to DLQ",
					"dlt_topic", msg.Topic+".DLT",
					"error", pubErr,
					"original_error", err,
				)
				return fmt.Errorf("dlq publish failed: %w (original: %v)", pubErr, err)
			}

			log.Warn("message sent to DLQ",
				"dlq_topic", msg.Topic+"DLQ",
				"attempts", attempts,
				"non_retryable", IsNonRetryable(err),
				"error", err,
			)

			// Сообщение надёжно лежит в DLQ — считаем обработанным, коммитим offset.
			return nil
		}
	}
}

// ---------- Recover ----------

func withRecover() Middleware {
	return func(next messaging.Handler) messaging.Handler {
		return func(ctx context.Context, msg messaging.Message) (err error) {
			defer func() {
				logger := logger.FromContext(ctx)
				if r := recover(); r != nil {
					logger.Error("handler panic",
						"panic", r, "stack", string(debug.Stack()))
					err = fmt.Errorf("panic: %v, %w", r, error_type.ErrPermanent)
				}
			}()
			return next(ctx, msg)
		}
	}
}
