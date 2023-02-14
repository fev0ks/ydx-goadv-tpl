package service

type OrderProcessingService interface {
	AddToProcessingQueue()
	ProcessOrder()
}
