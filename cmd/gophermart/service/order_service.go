package service

import (
	"context"
	"github.com/fev0ks/ydx-goadv-tpl/model"
	"github.com/fev0ks/ydx-goadv-tpl/repository"
	"github.com/theplant/luhn"
	"log"
)

type OrderService interface {
	SetOrder(ctx context.Context, userID int, orderID int) error
	ValidateOrder(ctx context.Context, orderID int) bool
	GetOrder(ctx context.Context, orderID int) (*model.Order, error)
	GetOrders(ctx context.Context, userID int) ([]*model.Order, error)
}

type orderService struct {
	orderRepo          repository.OrderRepository
	orderProcessingSrv OrderProcessingService
}

func NewOrderService(orderRepo repository.OrderRepository, orderProcessingSrv OrderProcessingService) OrderService {
	return &orderService{orderRepo, orderProcessingSrv}
}

func (os orderService) SetOrder(_ context.Context, userID int, orderID int) error {
	log.Printf("Additing order '%d' to accrual processing queue order, user: '%d'", orderID, userID)
	os.orderProcessingSrv.AddToAccrualProcessingQueue(userID, orderID)
	return nil
}

func (os orderService) ValidateOrder(_ context.Context, orderID int) bool {
	return luhn.Valid(orderID)
}

func (os orderService) GetOrder(ctx context.Context, orderID int) (*model.Order, error) {
	return os.orderRepo.GetOrder(ctx, orderID)
}

func (os orderService) GetOrders(ctx context.Context, userID int) ([]*model.Order, error) {
	return os.orderRepo.GetOrders(ctx, userID)
}
