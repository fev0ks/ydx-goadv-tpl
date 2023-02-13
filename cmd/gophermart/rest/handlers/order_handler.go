package handlers

import (
	"context"
	"fmt"
	"github.com/fev0ks/ydx-goadv-tpl/model/consts"
	"github.com/fev0ks/ydx-goadv-tpl/service"
	"io"
	"log"
	"net/http"
	"strconv"
)

type OrderHandler interface {
	SetOrderHandler() func(writer http.ResponseWriter, request *http.Request)
	GetOrdersHandler() func(writer http.ResponseWriter, request *http.Request)
}

type orderHandler struct {
	orderService service.OrderService
}

func NewOrderHandler(orderService service.OrderService) OrderHandler {
	return &orderHandler{orderService}
}

func (oh *orderHandler) SetOrderHandler() func(writer http.ResponseWriter, request *http.Request) {
	return func(writer http.ResponseWriter, request *http.Request) {
		ctx := request.Context()
		usernameCtx := ctx.Value(consts.UserCtxKey)
		if usernameCtx == nil {
			log.Printf("username is missed in context")
			http.Error(writer, "username is missed in context", http.StatusUnauthorized)
			return
		}
		username := usernameCtx.(string)
		order, err := oh.getOrder(ctx, request)
		if err != nil {
			log.Printf("failed to parse order request: %v", err)
			http.Error(writer, err.Error(), http.StatusBadRequest)
			return
		}

		err = oh.orderService.SetOrder(ctx, username, order)
		if err != nil {
			http.Error(writer, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

func (oh *orderHandler) getOrder(ctx context.Context, request *http.Request) (int, error) {
	body, err := io.ReadAll(request.Body)
	defer request.Body.Close()
	if err != nil {
		return 0, fmt.Errorf("request order in wrong format: %v", err)
	}
	if len(body) == 0 {
		return 0, fmt.Errorf("request body is empty")
	}
	order, err := strconv.Atoi(string(body))
	if err != nil {
		return 0, fmt.Errorf("request order in wrong format: %v", err)
	}
	if !oh.orderService.ValidateOrder(ctx, order) {
		return 0, fmt.Errorf("request order is not in Luna format")
	}
	return order, nil
}

func (oh *orderHandler) GetOrdersHandler() func(writer http.ResponseWriter, request *http.Request) {
	return func(writer http.ResponseWriter, request *http.Request) {
		ctx := request.Context()
		usernameCtx := ctx.Value(consts.UserCtxKey)
		if usernameCtx == nil {
			log.Printf("username is missed in context")
			http.Error(writer, "username is missed in context", http.StatusUnauthorized)
			return
		}
		username := usernameCtx.(string)
		orders, err := oh.orderService.GetOrders(ctx, username)
		if err != nil {
			log.Printf("failed to get orders for %s: %v", username, err)
			http.Error(writer, err.Error(), http.StatusInternalServerError)
			return
		}
		_, err = writer.Write([]byte(fmt.Sprintf("%v", orders)))
		if err != nil {
			http.Error(writer, err.Error(), http.StatusInternalServerError)
			return
		}
		writer.WriteHeader(http.StatusOK)
	}
}
