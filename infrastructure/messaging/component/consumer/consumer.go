package consumer

import (
	"context"
	"errors"
	"fmt"
	logger "go_project/infrastructure/logging"
	"go_project/infrastructure/messaging"
	"go_project/infrastructure/messaging/component/producer"
	"go_project/infrastructure/messaging/middleware"
	"io"
	"time"
)

type Reader interface {
	FetchMessage(ctx context.Context) (messaging.Delivery, error)
	Close() error
}

type Consumer struct {
	sub Reader
	h   messaging.Handler
}

func NewConsumer(
	sub Reader,
	pub producer.Writer,
	h messaging.Handler,
	log logger.Logger,
) *Consumer {

	return &Consumer{
		sub: sub,
		h:   middleware.Wrap(h, log, pub, "kafka")}
}

func (c *Consumer) Start(ctx context.Context) error {
	for {
		deliver, err := c.sub.FetchMessage(ctx)

		if err != nil {
			if ctx.Err() != nil || errors.Is(err, io.EOF) {
				return nil
			}
			return err
		}

		if err := c.h(ctx, deliver.Message); err != nil {
			return fmt.Errorf("handle %s: %w", deliver.Message.Topic, err)
		}

		cctx, cancel := context.WithTimeout(context.WithoutCancel(ctx), 5*time.Second)
		err = deliver.Commit(cctx)
		cancel()
		if err != nil {
			return err
		}
	}
}

func (c *Consumer) Stop(ctx context.Context) error { return c.sub.Close() }
