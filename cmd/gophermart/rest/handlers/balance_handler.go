package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/fev0ks/ydx-goadv-tpl/model"
	"github.com/fev0ks/ydx-goadv-tpl/model/consts"
	"github.com/fev0ks/ydx-goadv-tpl/model/consts/rest"
	"github.com/fev0ks/ydx-goadv-tpl/service"
	"log"
	"net/http"
)

type BalanceHandler interface {
	GetBalanceHandler() func(writer http.ResponseWriter, request *http.Request)
	BalanceWithdrawHandler() func(writer http.ResponseWriter, request *http.Request)
	GetWithdrawalsHandler() func(writer http.ResponseWriter, request *http.Request)
}

type balanceHandler struct {
	balanceService service.BalanceService
	orderService   service.OrderService
}

func NewBalanceHandler(balanceService service.BalanceService, orderService service.OrderService) BalanceHandler {
	return &balanceHandler{balanceService, orderService}
}

func (bh *balanceHandler) GetBalanceHandler() func(writer http.ResponseWriter, request *http.Request) {
	return func(writer http.ResponseWriter, request *http.Request) {
		ctx := request.Context()
		usernameCtx := ctx.Value(consts.UserIDCtxKey)
		if usernameCtx == nil {
			log.Printf("user is missed in context")
			http.Error(writer, "Unauthorized", http.StatusUnauthorized)
			return
		}
		userID := usernameCtx.(int)
		balance, err := bh.balanceService.GetBalance(ctx, userID)
		if err != nil {
			log.Printf("failed to get balance for %d: %v", userID, err)
			http.Error(writer, err.Error(), http.StatusInternalServerError)
			return
		}
		writer.Header().Add(rest.ContentType, rest.ApplicationJSON)
		err = json.NewEncoder(writer).Encode(balance)
		if err != nil {
			log.Printf("failed to write balance to response for %d: %v", userID, err)
			http.Error(writer, err.Error(), http.StatusInternalServerError)
			return
		}
		writer.WriteHeader(http.StatusOK)
	}
}

func (bh *balanceHandler) BalanceWithdrawHandler() func(writer http.ResponseWriter, request *http.Request) {
	return func(writer http.ResponseWriter, request *http.Request) {
		ctx := request.Context()
		usernameCtx := ctx.Value(consts.UserIDCtxKey)
		if usernameCtx == nil {
			log.Printf("user is missed in context")
			http.Error(writer, "Unauthorized", http.StatusUnauthorized)
			return
		}
		userID := usernameCtx.(int)
		withdrawRequest := &model.WithdrawRequest{}
		err := json.NewDecoder(request.Body).Decode(withdrawRequest)
		if err != nil {
			log.Printf("failed to parse withdraw request: %v", err)
			http.Error(writer, fmt.Sprintf("failed to parse withdraw request: %v", err), http.StatusBadRequest)
			return
		}
		if !bh.orderService.ValidateOrder(ctx, withdrawRequest.Order) {
			log.Printf("request '%d' order is not in Luna format", withdrawRequest.Order)
			http.Error(writer, fmt.Sprintf("request '%d' order is not in Luna format", withdrawRequest.Order), http.StatusUnprocessableEntity)
			return
		}
		order, err := bh.orderService.GetOrder(ctx, withdrawRequest.Order)
		if err != nil {
			http.Error(writer,
				fmt.Sprintf("failed to check '%d' order existance: %v", withdrawRequest.Order, err),
				http.StatusInternalServerError,
			)
			return
		}
		if order != nil {
			if order.UserID != userID {
				log.Printf("request order '%d' is already used by another user", withdrawRequest.Order)
				http.Error(writer, fmt.Sprintf("request order '%d' is already used by another user", withdrawRequest.Order), http.StatusConflict)
				return
			}
		}
		//ADD different errors
		err = bh.balanceService.BalanceWithdraw(ctx, userID, withdrawRequest)
		if err != nil {
			log.Printf("failed to parse withdraw request: %v", err)
			http.Error(writer, fmt.Sprintf("failed to parse withdraw request: %v", err), http.StatusBadRequest)
			return
		}
	}
}

func (bh *balanceHandler) GetWithdrawalsHandler() func(writer http.ResponseWriter, request *http.Request) {
	return func(writer http.ResponseWriter, request *http.Request) {
		ctx := request.Context()
		usernameCtx := ctx.Value(consts.UserIDCtxKey)
		if usernameCtx == nil {
			log.Printf("username is missed in context")
			http.Error(writer, "username is missed in context", http.StatusUnauthorized)
			return
		}
		userID := usernameCtx.(int)
		withdrawals, err := bh.balanceService.GetWithdrawals(ctx, userID)
		if err != nil {
			log.Printf("failed to get withdrawals for %d: %v", userID, err)
			http.Error(writer, err.Error(), http.StatusInternalServerError)
			return
		}
		if len(withdrawals) == 0 {
			writer.WriteHeader(http.StatusNoContent)
			return
		}
		writer.Header().Add(rest.ContentType, rest.ApplicationJSON)
		err = json.NewEncoder(writer).Encode(withdrawals)
		if err != nil {
			log.Printf("failed to write withdrawals to response for %d: %v", userID, err)
			http.Error(writer, err.Error(), http.StatusInternalServerError)
			return
		}
		writer.WriteHeader(http.StatusOK)

	}
}
