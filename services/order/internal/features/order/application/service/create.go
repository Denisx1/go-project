package service

import (
	"context"
	"fmt"
	"go_project/order/internal/core/domain"
)

func (s *Service) CreateOrder(
	ctx context.Context,
	order domain.Order,
) (domain.Order, error) {
	var resultOrder domain.Order
	err := s.TxManager.WithTx(ctx, func(ctx context.Context) error {
		newOrder, err := s.Repository.CreateOrder(ctx, order)
		if err != nil {
			return fmt.Errorf("create order %w", err)
		}
		resultOrder = newOrder
		if err := s.OutboxRepo.Publish(ctx, domain.NewOrderCreated(newOrder)); err != nil {
			return fmt.Errorf("create outbox: %w", err)
		}
		if err := s.Repository.DeleteOrder(ctx, newOrder.ID); err != nil {
			return fmt.Errorf("delete order: %w", err)
		}

		resultOrder = newOrder
		return nil
	})
	if err != nil {
		return domain.Order{}, err
	}
	fmt.Println(resultOrder)
	return resultOrder, nil
}
