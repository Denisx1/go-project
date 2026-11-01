package service

import (
	"context"
	"fmt"
)

func (s *Service) DeleteOrder(ctx context.Context, id string) error {
	if err := s.Repository.DeleteOrder(ctx, id); err != nil {
		return fmt.Errorf("delete order: %w", err)
	}
	return nil
}
