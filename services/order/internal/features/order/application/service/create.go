package service

import (
	"context"
	"fmt"
	"go_project/order/internal/core/domain"
)

func (s *Service) CreateOrder(
	ctx context.Context,
	order domain.Order,
) error {
	newOrder, err := s.Repository.CreateOrder(ctx, order)
	if err != nil {
		return fmt.Errorf("failed to create order: %w", err)
	}
	fmt.Println(newOrder)
	return nil
}
