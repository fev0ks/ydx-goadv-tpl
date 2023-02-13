package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/fev0ks/ydx-goadv-tpl/model"
	"github.com/fev0ks/ydx-goadv-tpl/service"
	"log"
	"net/http"
)

type UserHandler interface {
	RegisterHandler() func(writer http.ResponseWriter, request *http.Request)
	LoginHandler() func(writer http.ResponseWriter, request *http.Request)
}

type userHandler struct {
	sessionService service.SessionService
	userService    service.UserService
}

func NewUserHandler(
	sessionService service.SessionService,
	userService service.UserService,
) UserHandler {
	return &userHandler{sessionService, userService}
}

func (uh *userHandler) RegisterHandler() func(writer http.ResponseWriter, request *http.Request) {
	return func(writer http.ResponseWriter, request *http.Request) {
		ctx := request.Context()
		userRequest := &model.UserRequest{}
		err := json.NewDecoder(request.Body).Decode(userRequest)
		if err != nil {
			log.Printf("failed to parse user request: %v", err)
			writer.WriteHeader(http.StatusBadRequest)
			return
		}
		ok, err := uh.userService.IsExistUsername(ctx, userRequest.Username)
		if err != nil {
			log.Printf("failed to check user existance: %v", err)
			http.Error(writer, err.Error(), http.StatusInternalServerError)
			return
		}
		if ok {
			log.Printf("user '%s' already exist", userRequest.Username)
			http.Error(writer, fmt.Sprintf("user '%s' already exist", userRequest.Username), http.StatusConflict)
			return
		}
		err = uh.userService.CreateUser(ctx, userRequest)
		if err != nil {
			log.Printf("failed to create user: %v", err)
			http.Error(writer, err.Error(), http.StatusInternalServerError)
			return
		}

		sessionToken, expiresAt := uh.sessionService.CreateSession(userRequest.Username)
		http.SetCookie(writer, &http.Cookie{
			Name:     "session_token",
			Value:    sessionToken,
			Expires:  expiresAt,
			HttpOnly: true,
		})
	}
}

func (uh *userHandler) LoginHandler() func(writer http.ResponseWriter, request *http.Request) {
	return func(writer http.ResponseWriter, request *http.Request) {
		ctx := request.Context()
		userRequest := &model.UserRequest{}
		err := json.NewDecoder(request.Body).Decode(userRequest)
		if err != nil {
			log.Printf("failed to parse user request: %v", err)
			writer.WriteHeader(http.StatusBadRequest)
			return
		}
		ok, err := uh.userService.IsCorrectUserPassword(ctx, userRequest)
		if err != nil {
			log.Printf("failed to check user password: %v", err)
			http.Error(writer, err.Error(), http.StatusInternalServerError)
			return
		}
		if !ok {
			log.Printf("password is incorrect")
			writer.WriteHeader(http.StatusUnauthorized)
			return
		}
		sessionToken, expiresAt := uh.sessionService.CreateSession(userRequest.Username)
		http.SetCookie(writer, &http.Cookie{
			Name:     "session_token",
			Value:    sessionToken,
			Expires:  expiresAt,
			HttpOnly: true,
		})
	}
}
