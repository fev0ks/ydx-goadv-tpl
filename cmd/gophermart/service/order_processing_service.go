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
	AddToProcessingQueue(userOrder *model.UserOrder)
}

type orderProcessingService struct {
	mutex         *sync.RWMutex
	accrualClient clients.AccrualClient
	orderRepo     repository.OrderRepository
	queueCh       chan *model.UserOrder
	queue         []*model.UserOrder
}

func NewOrderProcessingService(
	ctx context.Context,
	accrualClient clients.AccrualClient,
	orderRepo repository.OrderRepository,
) OrderProcessingService {
	op := &orderProcessingService{
		&sync.RWMutex{},
		accrualClient,
		orderRepo,
		make(chan *model.UserOrder, 100),
		make([]*model.UserOrder, 0),
	}
	go op.orderQueueProcessing(ctx)
	//go op.queuePoller()
	return op
}

func (op *orderProcessingService) AddToProcessingQueue(userOrder *model.UserOrder) {
	//op.mutex.Lock()
	//op.queue = append(op.queue, userOrder)
	//op.mutex.Unlock()
	op.queueCh <- userOrder
}

//
////TODO bull shit
//func (op *orderProcessingService) queuePoller() {
//	for {
//		op.mutex.Lock()
//		defer op.mutex.Unlock()
//		if len(op.queue) > 0 {
//			op.queueCh <- op.queue[0]
//		}
//	}
//}

func (op *orderProcessingService) orderQueueProcessing(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case userOrder := <-op.queueCh:
			accrualOrder, err := op.accrualClient.GetOrderStatus(ctx, userOrder.Order)
			if err != nil {
				log.Printf("%v", err)
				retryError, ok := err.(*model.RetryError)
				if !ok {
					log.Printf("%v", err)
					continue
				}
				op.queueCh <- userOrder //lock?????
				time.Sleep(time.Second * time.Duration(retryError.RetryAfter))
			}
			err = op.saveOrder(ctx, accrualOrder)
			if err != nil {
				log.Printf("%v", err)
				continue
			}
			if accrualOrder.Status != model.ProcessedStatus && accrualOrder.Status != model.InvalidStatus {
				op.queueCh <- userOrder
			}
		}
	}
}

func (op *orderProcessingService) saveOrder(ctx context.Context, accrualOrder *model.AccrualOrder) error {
	order := &model.Order{
		Number:     accrualOrder.Order,
		Status:     accrualOrder.Status,
		Accrual:    accrualOrder.Accrual,
		UploadedAt: time.Now().Format(time.RFC3339),
	}
	return op.orderRepo.InsertOrUpdateOrder(ctx, "", order)
}
