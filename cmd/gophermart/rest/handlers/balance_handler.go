package handlers

import "net/http"

type BalanceHandler interface {
}

type balanceHandler struct {
}

func NewBalanceHandler() BalanceHandler {
	return &balanceHandler{}
}

func (bh *balanceHandler) GetBalanceHandler() func(writer http.ResponseWriter, request *http.Request) {
	return func(writer http.ResponseWriter, request *http.Request) {

	}
}

func (bh *balanceHandler) BalanceWithdrawHandler() func(writer http.ResponseWriter, request *http.Request) {
	return func(writer http.ResponseWriter, request *http.Request) {

	}
}

func (bh *balanceHandler) GetWithdrawsHistoryHandler() func(writer http.ResponseWriter, request *http.Request) {
	return func(writer http.ResponseWriter, request *http.Request) {

	}
}
