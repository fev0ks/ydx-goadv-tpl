package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/fev0ks/ydx-goadv-tpl/model/consts"
	"github.com/fev0ks/ydx-goadv-tpl/model/consts/rest"
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
		usernameCtx := ctx.Value(consts.UserIDCtxKey)
		if usernameCtx == nil {
			log.Printf("userID is missed in context")
			http.Error(writer, "userID is missed in context", http.StatusUnauthorized)
			return
		}
		userID := usernameCtx.(int)
		orderID, err := oh.getRequestOrder(request)
		if err != nil {
			log.Printf("failed to parse order request: %v", err)
			http.Error(writer, fmt.Sprintf("failed to parse order request: %v", err), http.StatusBadRequest)
			return
		}
		if !oh.orderService.ValidateOrder(ctx, orderID) {
			log.Printf("request '%d' order is not in Luna format", orderID)
			http.Error(writer, fmt.Sprintf("request '%d' order is not in Luna format", orderID), http.StatusUnprocessableEntity)
			return
		}
		order, err := oh.orderService.GetOrder(ctx, orderID)
		if err != nil {
			http.Error(writer,
				fmt.Sprintf("failed to check '%d' order existance: %v", orderID, err),
				http.StatusInternalServerError,
			)
			return
		}
		if order != nil {
			if order.UserID != userID {
				log.Printf("request order '%d' is already used by another user", orderID)
				http.Error(writer, fmt.Sprintf("request order '%d' is already used by another user", orderID), http.StatusConflict)
				return
			} else {
				log.Printf("request '%d' order is already used by current user", orderID)
				http.Error(writer, fmt.Sprintf("request '%d' order is already used by current user", orderID), http.StatusOK)
				return
			}
		}
		_ = oh.orderService.SetOrder(ctx, userID, orderID)
		writer.WriteHeader(http.StatusAccepted)
	}
}

func (oh *orderHandler) getRequestOrder(request *http.Request) (int, error) {
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
	return order, nil
}

func (oh *orderHandler) GetOrdersHandler() func(writer http.ResponseWriter, request *http.Request) {
	return func(writer http.ResponseWriter, request *http.Request) {
		ctx := request.Context()
		usernameCtx := ctx.Value(consts.UserIDCtxKey)
		if usernameCtx == nil {
			log.Printf("username is missed in context")
			http.Error(writer, "username is missed in context", http.StatusUnauthorized)
			return
		}
		userID := usernameCtx.(int)
		orders, err := oh.orderService.GetOrders(ctx, userID)
		if err != nil {
			log.Printf("failed to get orders for %d: %v", userID, err)
			http.Error(writer, err.Error(), http.StatusInternalServerError)
			return
		}
		if len(orders) == 0 {
			writer.WriteHeader(http.StatusNoContent)
			return
		}
		writer.Header().Add(rest.ContentType, rest.ApplicationJSON)
		err = json.NewEncoder(writer).Encode(orders)
		if err != nil {
			log.Printf("failed to write orders to response for %d: %v", userID, err)
			http.Error(writer, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}
