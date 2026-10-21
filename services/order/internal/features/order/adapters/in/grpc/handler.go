package OrderGrpcHandler

import orderv1 "go_project/contracts/gen/go/order"

type OrderGRPCHandler struct {
	orderv1.UnimplementedOrderServer
}

func NewOrderGRPCHandler() *OrderGRPCHandler {
	return &OrderGRPCHandler{}
}
