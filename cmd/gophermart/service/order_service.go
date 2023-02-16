package service

import (
	"context"
	"fmt"
	"github.com/fev0ks/ydx-goadv-tpl/model"
	"github.com/fev0ks/ydx-goadv-tpl/repository"
	"github.com/fev0ks/ydx-goadv-tpl/rest/clients"
	"github.com/theplant/luhn"
	"time"
)

type OrderService interface {
	SetOrder(ctx context.Context, userId int, orderId int) error
	ValidateOrder(ctx context.Context, orderId int) bool
	IsOrderUsed(ctx context.Context, orderId int) (bool, error)
	GetOrders(ctx context.Context, userId int) ([]*model.Order, error)
}

type orderService struct {
	orderRepo     repository.OrderRepository
	accrualClient clients.AccrualClient
}

func NewOrderService(orderRepo repository.OrderRepository, accrualClient clients.AccrualClient) OrderService {
	return &orderService{orderRepo, accrualClient}
}

func (os orderService) SetOrder(ctx context.Context, userId int, orderId int) error {
	accrualOrder, err := os.accrualClient.GetOrderStatus(ctx, fmt.Sprintf("%d", orderId))
	if err != nil {
		return err
	}
	order := &model.Order{
		Number:     accrualOrder.Order,
		Status:     accrualOrder.Status,
		Accrual:    accrualOrder.Accrual,
		UploadedAt: time.Now(),
	}
	return os.orderRepo.InsertOrder(ctx, userId, order)
}

func (os orderService) ValidateOrder(_ context.Context, orderId int) bool {
	return luhn.Valid(orderId)
}

func (os orderService) IsOrderUsed(ctx context.Context, orderId int) (bool, error) {
	return os.orderRepo.IsOrderExist(ctx, orderId)
}

func (os orderService) GetOrders(ctx context.Context, userId int) ([]*model.Order, error) {
	return os.orderRepo.GetOrders(ctx, userId)
}
