package domain

type Event interface {
}

type OrderCreated struct {
	Order Order
}

func NewOrderCreated(order Order) OrderCreated {
	return OrderCreated{
		Order: order,
	}
}
