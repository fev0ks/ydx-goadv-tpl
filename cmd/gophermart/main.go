package main

import (
	"context"
	"github.com/fev0ks/ydx-goadv-tpl/config"
	"github.com/fev0ks/ydx-goadv-tpl/repository"
	"github.com/fev0ks/ydx-goadv-tpl/rest"
	"github.com/fev0ks/ydx-goadv-tpl/rest/clients"
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
	ctx := context.Background()
	log.Printf("Server args: %s", os.Args[1:])
	appConfig := config.InitAppConfig()

	DBProvider, err := repository.NewPgProvider(ctx, appConfig)
	if err != nil {
		log.Fatalln(err)
	}
	sessionStorage := storage.NewSessionStorage()
	userRepo := repository.NewUserRepository(DBProvider)
	orderRepo := repository.NewOrderRepository(DBProvider)

	accrualClient := clients.NewAccrualClient(clients.CreateClient(appConfig.AccrualAddress))
	sessionService := service.NewSessionService(sessionStorage, appConfig.SessionLifetime)

	userService := service.NewUserService(userRepo)
	orderService := service.NewOrderService(orderRepo, accrualClient)

	router := rest.NewRouter()

	userHandler := handlers.NewUserHandler(sessionService, userService)
	orderHandler := handlers.NewOrderHandler(orderService)
	balanceHandler := handlers.NewBalanceHandler()
	healthChecker := rest.NewHealthChecker(ctx, DBProvider)

	authMiddleware := middlewares.NewAuthMiddleware(sessionService)
	rest.HandleUserRequests(router, userHandler)
	rest.HandleOrderRequests(router, authMiddleware, orderHandler)
	rest.HandleBalanceRequests(router, authMiddleware, balanceHandler)
	rest.HandleHeathCheck(router, healthChecker)

	shutdown.ProperExitDefer(&shutdown.ExitHandler{
		//ToStop:    stopCh,
		//ToExecute: toExecute,
		//ToClose:   []io.Closer{repository},
	})
	log.Printf("Server started on %s", appConfig.ServerAddress)
	log.Fatal(http.ListenAndServe(appConfig.ServerAddress, router))
}
