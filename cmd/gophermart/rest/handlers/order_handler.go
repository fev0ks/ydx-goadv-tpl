package handlers

import "net/http"

type OrderHandler interface {
}

type orderHandler struct {
}

func NewOrderHandler() OrderHandler {
	return &orderHandler{}
}

func (oh *orderHandler) CreateOrderHandler() func(writer http.ResponseWriter, request *http.Request) {
	return func(writer http.ResponseWriter, request *http.Request) {

	}
}

func (oh *orderHandler) GetOrderHandler() func(writer http.ResponseWriter, request *http.Request) {
	return func(writer http.ResponseWriter, request *http.Request) {

	}
}
