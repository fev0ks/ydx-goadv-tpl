package handlers

import "net/http"

type UserHandler interface {
}

type userHandler struct {
}

func NewUserHandler() UserHandler {
	return &userHandler{}
}

func (uh *userHandler) RegisterHandler() func(writer http.ResponseWriter, request *http.Request) {
	return func(writer http.ResponseWriter, request *http.Request) {

	}
}

func (uh *userHandler) LoginHandler() func(writer http.ResponseWriter, request *http.Request) {
	return func(writer http.ResponseWriter, request *http.Request) {

	}
}
