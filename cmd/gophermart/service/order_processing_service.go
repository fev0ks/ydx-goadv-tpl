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

// OrderProcessingService TODO looks very bad
type OrderProcessingService interface {
	AddToAccrualOrderProcessingQueue(order *model.Order)
}

type orderProcessingService struct {
	*sync.RWMutex
	queue         []*model.Order
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
		make([]*model.Order, 0),
		accrualClient,
		orderRepo,
	}
	go op.orderQueueProcessing(ctx)
	return op
}

func (op *orderProcessingService) AddToAccrualOrderProcessingQueue(order *model.Order) {
	op.Lock()
	op.queue = append(op.queue, order)
	op.Unlock()
}

func (op *orderProcessingService) orderQueueProcessing(ctx context.Context) {
	ticker := time.NewTicker(time.Second)
	for {
		select {
		case <-ctx.Done():
			log.Println("Order Queue processing is canceled")
			return
		case <-ticker.C:
			toRetry := make([]*model.Order, 0)
			for _, order := range op.retrieveQueueData() {
				accrualOrder, err := op.accrualClient.GetOrderStatus(ctx, order.Number)
				if err != nil {
					retryError, ok := err.(*model.RetryError)
					if !ok {
						log.Printf("failed to get order accrual status: %v", err)
						continue
					}
					log.Printf("Accrual warning StatusTooManyRequests: %v", err)
					toRetry = append(toRetry, order)
					waitingDuration := time.Second * time.Duration(retryError.RetryAfter)
					log.Printf("Waiting for: %v", waitingDuration)
					time.Sleep(waitingDuration)
					continue
				}
				if order.Status != accrualOrder.Status {
					err = op.updateOrderState(ctx, order, accrualOrder)
				}
				if err != nil {
					log.Printf("failed to save order %v: %v", accrualOrder, err)
					continue
				}
				if accrualOrder.Status != model.InvalidStatus && accrualOrder.Status != model.ProcessedStatus {
					toRetry = append(toRetry, order)
				}
			}
			op.backToQueue(toRetry)
		}
	}
}

func (op *orderProcessingService) retrieveQueueData() []*model.Order {
	op.Lock()
	defer op.Unlock()
	data := op.queue
	op.queue = op.queue[0:0]
	return data
}

func (op *orderProcessingService) backToQueue(backToQueue []*model.Order) {
	op.Lock()
	defer op.Unlock()
	op.queue = append(op.queue, backToQueue...)
}

func (op *orderProcessingService) updateOrderState(
	ctx context.Context,
	order *model.Order,
	accrualOrder *model.AccrualOrder,
) error {
	order.Status = accrualOrder.Status
	order.Accrual = accrualOrder.Accrual
	return op.orderRepo.UpdateOrder(ctx, order)
}
