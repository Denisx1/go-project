package service

import (
	"context"
	"go_project/order/internal/core/domain"
)

type TxManager interface {
	WithTx(ctx context.Context, fn func(ctx context.Context) error) error
	WithTxSerializable(ctx context.Context, fn func(ctx context.Context) error) error
}

type Repository interface {
	CreateOrder(ctx context.Context, order domain.Order) (domain.Order, error)
	DeleteOrder(ctx context.Context, orderId string) error
}
type OutboxRepository interface {
	Publish(
		ctx context.Context,
		e domain.Event,
	) error
}

type Service struct {
	Repository Repository
	OutboxRepo OutboxRepository
	TxManager  TxManager
}

func NewService(
	repository Repository,
	outboxRepo OutboxRepository,
	txManager TxManager,
) *Service {
	return &Service{
		Repository: repository,
		OutboxRepo: outboxRepo,
		TxManager:  txManager,
	}
}
