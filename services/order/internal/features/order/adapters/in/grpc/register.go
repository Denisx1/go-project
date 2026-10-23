package OrderGrpcHandler

import (
	"context"
	orderv1 "go_project/contracts/gen/go/order"
	GrpcResponse "go_project/infrastructure/grpc/server/response"
	logger "go_project/infrastructure/logging"
	"go_project/order/internal/core/domain"
)

func (h *OrderGRPCHandler) Register(
	ctx context.Context,
	req *orderv1.RegisterRequest,
) (*orderv1.RegisterResponse, error) {
	log := logger.FromContext(ctx)
	responseHandler := GrpcResponse.NewGrpcResponseHandler(log)
	orderDomain := domain.NewOrderUninitialized(req.Title, req.Text)
	created, err := h.OrderService.CreateOrder(ctx, orderDomain)
	if err != nil {
		return nil, responseHandler.ErrorResponse(err, "failed to create order")
	}
	return &orderv1.RegisterResponse{
		Id:       created.ID,
		Version:  created.Version,
		Title:    created.Title,
		Text:     created.Text,
		IsActive: created.IsActive,
	}, nil

}
