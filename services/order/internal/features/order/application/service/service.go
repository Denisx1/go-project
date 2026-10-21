package service

import (
	"context"
	"go_project/order/internal/core/domain"
)

type Repository interface {
	CreateOrder(ctx context.Context, order domain.Order) (domain.Order, error)
}

type Service struct {
	Repository Repository
}

func NewService(repository Repository) *Service {
	return &Service{
		Repository: repository,
	}
}
