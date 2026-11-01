package grpchandler

import (
	"context"
	"fmt"
	orderv1 "go_project/contracts/gen/go/order"
)

func (h *OrderGRPCHandler) Delete(
	ctx context.Context,
	req *orderv1.DeleteRequest,
) (*orderv1.DeleteResponse, error) {
	if err := h.OrderService.DeleteOrder(ctx, req.Id); err != nil {
		return nil, fmt.Errorf("failed to delete order: %w", err)
	}
	return &orderv1.DeleteResponse{}, nil

}
