package grpchandler

import (
	"context"
	"fmt"
	orderv1 "go_project/contracts/gen/go/order"
)

func (h *OrderGRPCHandler) Get(
	ctx context.Context,
	req *orderv1.GetRequest,
) (*orderv1.GetResponse, error) {

	order, err := h.OrderService.GetOrder(ctx, req.Id)
	if err != nil {
		return nil, fmt.Errorf("get order: %w", err)
	}
	return &orderv1.GetResponse{
		Id:       order.ID,
		Version:  order.Version,
		Title:    order.Title,
		Text:     order.Text,
		IsActive: order.IsActive,
	}, nil

}
