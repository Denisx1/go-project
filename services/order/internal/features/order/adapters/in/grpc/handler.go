package OrderGrpcHandler

import (
	"context"
	orderv1 "go_project/contracts/gen/go/order"
	"go_project/order/internal/core/domain"
)

type OrderService interface {
	CreateOrder(
		ctx context.Context,
		order domain.Order,
	) (domain.Order, error)
}
type OrderGRPCHandler struct {
	OrderService OrderService
	orderv1.UnimplementedOrderServer
}

func NewOrderGRPCHandler(orderService OrderService) *OrderGRPCHandler {
	return &OrderGRPCHandler{
		OrderService: orderService,
	}
}
