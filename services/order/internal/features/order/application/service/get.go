package service

import (
	"context"
	"fmt"
	"go_project/order/internal/core/domain"
)

func (s *Service) GetOrder(
	ctx context.Context, id string,
) (domain.Order, error) {

	order, err := s.Repository.GetOrder(ctx, id)
	if err != nil {
		return domain.Order{}, fmt.Errorf("get order from repository: %w", err)
	}
	return order, nil
}
