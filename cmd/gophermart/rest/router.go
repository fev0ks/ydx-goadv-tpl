package rest

import (
	"github.com/fev0ks/ydx-goadv-tpl/model/consts/rest"
	"github.com/fev0ks/ydx-goadv-tpl/rest/handlers"
	"github.com/fev0ks/ydx-goadv-tpl/rest/middlewares"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func NewRouter() chi.Router {
	router := chi.NewRouter()
	router.Use(middlewares.TimerTrace)
	router.Use(middleware.RealIP)
	router.Use(middleware.Recoverer)
	router.Use(middleware.Compress(3, rest.ApplicationJSON, rest.TextPlain))
	router.Use(middlewares.Decompress)
	return router
}

func HandleUserRequests(router chi.Router, userService handlers.UserHandler) {
	router.Route("/api/user", func(r chi.Router) {
		//router.Post("/register", userService.ReceptionMetricHandler())
		//router.Post("/login", userService.ReceptionMetricHandler())
	})
}

func HandleOrderRequests(router chi.Router, orderHandler handlers.OrderHandler) {
	router.Route("/api/user/orders", func(r chi.Router) {
		//router.Post("/", userService.ReceptionMetricHandler())
		//router.Get("/", userService.ReceptionMetricHandler())
	})
}

func HandleBalanceRequests(router chi.Router, balanceHandler handlers.BalanceHandler) {
	router.Route("/api/user/balance", func(r chi.Router) {
		//router.Get("/balance", userService.ReceptionMetricHandler())
		//router.Post("/balance/withdraw", userService.ReceptionMetricHandler())
	})
	router.Route("/api/user/withdraws", func(r chi.Router) {
		//router.Get("/", userService.ReceptionMetricHandler())
	})
}

func HandleHeathCheck(router chi.Router, hc HealthChecker) {
	router.Get("/ping", hc.CheckDBHandler())
}
