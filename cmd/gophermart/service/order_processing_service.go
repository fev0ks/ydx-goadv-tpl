package service

import (
	"context"
	"github.com/fev0ks/ydx-goadv-tpl/model"
	"github.com/fev0ks/ydx-goadv-tpl/repository"
	"github.com/fev0ks/ydx-goadv-tpl/rest/clients"
	"log"
	"sync"
	"time"
)

type OrderProcessingService interface {
	AddToAccrualProcessingQueue(userID, orderID int)
}

type orderProcessingService struct {
	*sync.RWMutex
	queue         []*model.UserOrder
	accrualClient clients.AccrualClient
	orderRepo     repository.OrderRepository
}

func NewOrderProcessingService(
	ctx context.Context,
	accrualClient clients.AccrualClient,
	orderRepo repository.OrderRepository,
) OrderProcessingService {
	op := &orderProcessingService{
		&sync.RWMutex{},
		make([]*model.UserOrder, 0),
		accrualClient,
		orderRepo,
	}
	go op.orderQueueProcessing(ctx)
	return op
}

func (op *orderProcessingService) AddToAccrualProcessingQueue(userID, orderID int) {
	userOrder := &model.UserOrder{
		UserID: userID,
		Order: &model.Order{
			Number:     orderID,
			UploadedAt: time.Now(),
		},
	}
	op.Lock()
	op.queue = append(op.queue, userOrder)
	op.Unlock()
}

func (op *orderProcessingService) orderQueueProcessing(ctx context.Context) {
	ticker := time.NewTicker(time.Second / 2)
	for {
		select {
		case <-ctx.Done():
			log.Println("Order Queue processing is canceled")
			return
		case <-ticker.C:
			toRetry := make([]*model.UserOrder, 0)
			for _, userOrder := range op.retrieveQueueData() {
				accrualOrder, err := op.accrualClient.GetOrderStatus(ctx, userOrder.Order.Number)
				if err != nil {
					retryError, ok := err.(*model.RetryError)
					if !ok {
						log.Printf("failed to get order accrual status: %v", err)
						continue
					}
					log.Printf("Accrual warning StatusTooManyRequests: %v", err)
					toRetry = append(toRetry, userOrder)
					waitingDuration := time.Second * time.Duration(retryError.RetryAfter)
					log.Printf("Waiting for: %v", waitingDuration)
					time.Sleep(waitingDuration)
					continue
				}
				if userOrder.Order.Status == "" {
					err = op.saveOrder(ctx, userOrder.UserID, accrualOrder)
				} else if userOrder.Order.Status != accrualOrder.Status {
					err = op.updateOrder(ctx, userOrder.UserID, accrualOrder)
				}
				if err != nil {
					log.Printf("failed to save order %v: %v", accrualOrder, err)
					continue
				}
				if accrualOrder.Status != model.InvalidStatus && accrualOrder.Status != model.ProcessedStatus {
					toRetry = append(toRetry, userOrder)
				}
			}
			op.backToQueue(toRetry)
		}
	}
}

func (op *orderProcessingService) retrieveQueueData() []*model.UserOrder {
	op.Lock()
	defer op.Unlock()
	data := op.queue
	op.queue = op.queue[0:0]
	return data
}

func (op *orderProcessingService) backToQueue(backToQueue []*model.UserOrder) {
	op.Lock()
	defer op.Unlock()
	op.queue = append(op.queue, backToQueue...)
}

func (op *orderProcessingService) saveOrder(ctx context.Context, userID int, accrualOrder *model.AccrualOrder) error {
	order := accrualOrder.ToOrder()
	order.UploadedAt = time.Now()
	return op.orderRepo.InsertOrder(ctx, userID, order)
}

func (op *orderProcessingService) updateOrder(ctx context.Context, userID int, accrualOrder *model.AccrualOrder) error {
	return op.orderRepo.UpdateOrder(ctx, userID, accrualOrder.ToOrder())
}
