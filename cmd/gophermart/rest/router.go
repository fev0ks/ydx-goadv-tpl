package rest

import (
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
	//router.Use(middleware.Compress(3, rest.ApplicationJSON, rest.TextPlain))
	//router.Use(middlewares.Decompress)
	return router
}

func HandleUserRequests(
	router chi.Router,
	userHandler handlers.UserHandler,
) {
	router.Group(func(r chi.Router) {
		r.Route("/api/user", func(r chi.Router) {
			r.Post("/register", userHandler.RegisterHandler())
			r.Post("/login", userHandler.LoginHandler())
		})
	})
}

func HandleOrderRequests(router chi.Router,
	authMiddleware middlewares.SessionTokenValidator,
	orderHandler handlers.OrderHandler,
) {
	router.Group(func(r chi.Router) {
		r.Use(authMiddleware.ValidateSessionToken)
		r.Route("/api/user/orders", func(r chi.Router) {
			r.Post("/", orderHandler.SetOrderHandler())
			r.Get("/", orderHandler.GetOrdersHandler())
		})
	})
}

func HandleBalanceRequests(router chi.Router,
	authMiddleware middlewares.SessionTokenValidator,
	balanceHandler handlers.BalanceHandler,
) {
	router.Group(func(r chi.Router) {
		r.Use(authMiddleware.ValidateSessionToken)
		r.Route("/api/user/balance", func(r chi.Router) {
			r.Get("/", balanceHandler.GetBalanceHandler())
			r.Post("/withdraw", balanceHandler.BalanceWithdrawHandler())
		})
		r.Get("/api/user/withdrawals", balanceHandler.GetWithdrawalsHandler())
	})
}

func HandleHeathCheck(router chi.Router, hc HealthChecker) {
	router.Get("/ping", hc.CheckDBHandler())
}
