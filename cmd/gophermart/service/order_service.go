package service

import (
	"context"
	"github.com/fev0ks/ydx-goadv-tpl/repository"
	"github.com/theplant/luhn"
)

type OrderService interface {
	SetOrder(ctx context.Context, username string, order int) error
	ValidateOrder(ctx context.Context, order int) bool
	GetOrders(ctx context.Context, username string) ([]int, error)
}

type orderService struct {
	orderRepo repository.OrderRepository
}

func NewOrderService(orderRepo repository.OrderRepository) OrderService {
	return &orderService{orderRepo}
}

func (os orderService) SetOrder(ctx context.Context, username string, order int) error {
	return os.orderRepo.CreateOrder(ctx, username, order)
}

func (os orderService) ValidateOrder(_ context.Context, order int) bool {
	return luhn.Valid(order)
}

func (os orderService) GetOrders(ctx context.Context, username string) ([]int, error) {
	return os.orderRepo.GetOrders(ctx, username)
}
