package main

import (
	"context"
	"github.com/fev0ks/ydx-goadv-tpl/model/config"
	"github.com/fev0ks/ydx-goadv-tpl/repository"
	"github.com/fev0ks/ydx-goadv-tpl/rest"
	"github.com/fev0ks/ydx-goadv-tpl/rest/handlers"
	"github.com/fev0ks/ydx-goadv-tpl/shutdown"
	"log"
	"net/http"
	"os"
)

func main() {
	ctx := context.WithValue(context.Background(), "service", "tpl")
	//var err error
	log.Printf("Server args: %s", os.Args[1:])
	appConfig := config.InitAppConfig()

	dbProvider := repository.NewPgProvider()
	healthChecker := rest.NewHealthChecker(ctx, dbProvider)
	router := rest.NewRouter()
	userHandler := handlers.NewUserHandler()
	orderHandler := handlers.NewOrderHandler()
	balanceHandler := handlers.NewBalanceHandler()
	rest.HandleUserRequests(router, userHandler)
	rest.HandleOrderRequests(router, orderHandler)
	rest.HandleBalanceRequests(router, balanceHandler)
	rest.HandleHeathCheck(router, healthChecker)

	shutdown.ProperExitDefer(&shutdown.ExitHandler{
		//ToStop:    stopCh,
		//ToExecute: toExecute,
		//ToClose:   []io.Closer{repository},
	})
	log.Printf("Server started on %s", appConfig.ServerAddress)
	log.Fatal(http.ListenAndServe(appConfig.ServerAddress, router))
}
