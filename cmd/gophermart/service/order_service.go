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
	SetOrder(ctx context.Context, userID int, orderID int) error
	ValidateOrder(ctx context.Context, orderID int) bool
	IsOrderUsed(ctx context.Context, orderID int) (bool, error)
	GetOrders(ctx context.Context, userID int) ([]*model.Order, error)
}

type orderService struct {
	orderRepo     repository.OrderRepository
	accrualClient clients.AccrualClient
}

func NewOrderService(orderRepo repository.OrderRepository, accrualClient clients.AccrualClient) OrderService {
	return &orderService{orderRepo, accrualClient}
}

func (os orderService) SetOrder(ctx context.Context, userID int, orderID int) error {
	accrualOrder, err := os.accrualClient.GetOrderStatus(ctx, orderID)
	if err != nil {
		return err
	}
	order := &model.Order{
		Number:     accrualOrder.Order,
		Status:     accrualOrder.Status,
		Accrual:    accrualOrder.Accrual,
		UploadedAt: time.Now(),
	}
	return os.orderRepo.InsertOrder(ctx, userID, order)
}

func (os orderService) ValidateOrder(_ context.Context, orderID int) bool {
	return luhn.Valid(orderID)
}

func (os orderService) IsOrderUsed(ctx context.Context, orderID int) (bool, error) {
	return os.orderRepo.IsOrderExist(ctx, orderID)
}

func (os orderService) GetOrders(ctx context.Context, userID int) ([]*model.Order, error) {
	return os.orderRepo.GetOrders(ctx, userID)
}
