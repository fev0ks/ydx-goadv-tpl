package service

import (
	"context"
	"github.com/fev0ks/ydx-goadv-tpl/model"
	"github.com/fev0ks/ydx-goadv-tpl/repository"
	"github.com/fev0ks/ydx-goadv-tpl/rest/clients"
	"github.com/theplant/luhn"
	"time"
)

type OrderService interface {
	SetOrder(ctx context.Context, username string, orderNumber int) error
	ValidateOrder(ctx context.Context, orderNumber int) bool
	IsOrderUsed(ctx context.Context, orderNumber int) (bool, error)
	GetOrders(ctx context.Context, username string) ([]*model.Order, error)
}

type orderService struct {
	orderRepo     repository.OrderRepository
	accrualClient clients.AccrualClient
}

func NewOrderService(orderRepo repository.OrderRepository, accrualClient clients.AccrualClient) OrderService {
	return &orderService{orderRepo, accrualClient}
}

func (os orderService) SetOrder(ctx context.Context, username string, orderNumber int) error {
	accrualOrder, err := os.accrualClient.GetOrderStatus(ctx, orderNumber)
	if err != nil {
		return err
	}
	order := &model.Order{
		Number:     accrualOrder.Order,
		Status:     accrualOrder.Status,
		Accrual:    accrualOrder.Accrual,
		UploadedAt: time.Now().Format(time.RFC3339),
	}
	return os.orderRepo.InsertOrUpdateOrder(ctx, username, order)
}

func (os orderService) ValidateOrder(_ context.Context, orderNumber int) bool {
	return luhn.Valid(orderNumber)
}

func (os orderService) IsOrderUsed(ctx context.Context, orderNumber int) (bool, error) {
	return os.orderRepo.IsOrderExist(ctx, orderNumber)
}

func (os orderService) GetOrders(ctx context.Context, username string) ([]*model.Order, error) {
	return os.orderRepo.GetOrders(ctx, username)
}
