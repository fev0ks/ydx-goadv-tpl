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
			http.Error(writer, fmt.Sprintf("failed to parse user request: %v", err), http.StatusBadRequest)
			return
		}
		if userRequest.Login == "" || userRequest.Password == "" {
			log.Printf("login or password is empty")
			http.Error(writer, "login or password is empty", http.StatusBadRequest)
			return
		}
		ok, err := uh.userService.IsExistUsername(ctx, userRequest.Login)
		if err != nil {
			log.Printf("failed to check user existance: %v", err)
			http.Error(writer, fmt.Sprintf("failed to check user existance: %v", err), http.StatusInternalServerError)
			return
		}
		if ok {
			log.Printf("user '%s' already exist", userRequest.Login)
			http.Error(writer, fmt.Sprintf("user '%s' already exist", userRequest.Login), http.StatusConflict)
			return
		}
		user, err := uh.userService.CreateUser(ctx, userRequest)
		if err != nil {
			log.Printf("failed to create user: %v", err)
			http.Error(writer, fmt.Sprintf("failed to create user: %v", err), http.StatusInternalServerError)
			return
		}

		sessionToken, expiresAt := uh.sessionService.CreateSession(user.UserID)
		http.SetCookie(writer, &http.Cookie{
			Name:     "session_token",
			Value:    sessionToken,
			Expires:  expiresAt,
			HttpOnly: true,
		})
		log.Printf("'%s' was signed up", userRequest.Login)
	}
}

func (uh *userHandler) LoginHandler() func(writer http.ResponseWriter, request *http.Request) {
	return func(writer http.ResponseWriter, request *http.Request) {
		ctx := request.Context()
		userRequest := &model.UserRequest{}
		err := json.NewDecoder(request.Body).Decode(userRequest)
		if err != nil {
			log.Printf("failed to parse user request: %v", err)
			http.Error(writer, fmt.Sprintf("failed to parse user request: %v", err), http.StatusBadRequest)
			return
		}
		user, err := uh.userService.GetUser(ctx, userRequest.Login)
		if err != nil {
			log.Printf("failed to get user: %v", err)
			http.Error(writer, fmt.Sprintf("failed to get user: %v", err), http.StatusInternalServerError)
		}
		ok, err := uh.userService.ValidatePassword(ctx, user, userRequest.Password)
		if err != nil {
			log.Printf("failed to check user password: %v", err)
			http.Error(writer, fmt.Sprintf("failed to check user password: %v", err), http.StatusInternalServerError)
			return
		}
		if !ok {
			log.Printf("password is incorrect")
			http.Error(writer, "password is incorrect", http.StatusUnauthorized)
			return
		}
		sessionToken, expiresAt := uh.sessionService.CreateSession(user.UserID)
		http.SetCookie(writer, &http.Cookie{
			Name:     "session_token",
			Value:    sessionToken,
			Expires:  expiresAt,
			HttpOnly: true,
		})
		log.Printf("'%s' was signed in", userRequest.Login)
	}
}
