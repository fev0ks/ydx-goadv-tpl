package main

import (
	"context"
	"github.com/fev0ks/ydx-goadv-tpl/model/config"
	"github.com/fev0ks/ydx-goadv-tpl/repository"
	"github.com/fev0ks/ydx-goadv-tpl/rest"
	"github.com/fev0ks/ydx-goadv-tpl/rest/handlers"
	"github.com/fev0ks/ydx-goadv-tpl/rest/middlewares"
	"github.com/fev0ks/ydx-goadv-tpl/service"
	"github.com/fev0ks/ydx-goadv-tpl/shutdown"
	"github.com/fev0ks/ydx-goadv-tpl/storage"
	"log"
	"net/http"
	"os"
)

func main() {
	ctx := context.WithValue(context.Background(), "service", "tpl")
	log.Printf("Server args: %s", os.Args[1:])
	appConfig := config.InitAppConfig()

	dbProvider, err := repository.NewPgProvider(ctx, appConfig)
	if err != nil {
		log.Fatalln(err)
	}
	sessionStorage := storage.NewSessionStorage()
	userRepo := repository.NewUserRepository(dbProvider)
	sessionService := service.NewSessionService(sessionStorage, appConfig.SessionLifetime)
	userService := service.NewUserService(userRepo)

	router := rest.NewRouter()

	userHandler := handlers.NewUserHandler(sessionService, userService)
	orderHandler := handlers.NewOrderHandler()
	balanceHandler := handlers.NewBalanceHandler()
	healthChecker := rest.NewHealthChecker(ctx, dbProvider)

	tokenValidator := middlewares.NewAuthMiddleware(sessionService)
	rest.HandleUserRequests(router, userHandler)
	rest.HandleOrderRequests(router, tokenValidator, orderHandler)
	rest.HandleBalanceRequests(router, tokenValidator, balanceHandler)
	rest.HandleHeathCheck(router, healthChecker)

	shutdown.ProperExitDefer(&shutdown.ExitHandler{
		//ToStop:    stopCh,
		//ToExecute: toExecute,
		//ToClose:   []io.Closer{repository},
	})
	log.Printf("Server started on %s", appConfig.ServerAddress)
	log.Fatal(http.ListenAndServe(appConfig.ServerAddress, router))
}
