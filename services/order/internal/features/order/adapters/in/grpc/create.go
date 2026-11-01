package grpchandler

import (
	"context"
	"fmt"
	orderv1 "go_project/contracts/gen/go/order"

	"go_project/order/internal/core/domain"
)

func (h *OrderGRPCHandler) Create(
	ctx context.Context,
	req *orderv1.CreateRequest,
) (*orderv1.RegisterResponse, error) {
	orderDomain := domain.NewOrderUninitialized(req.Title, req.Text)
	created, err := h.OrderService.CreateOrder(ctx, orderDomain)
	if err != nil {
		return nil, fmt.Errorf("failed to create order: %w", err)
	}
	return &orderv1.RegisterResponse{
		Id:       created.ID,
		Version:  created.Version,
		Title:    created.Title,
		Text:     created.Text,
		IsActive: created.IsActive,
	}, nil

}
