package main

import (
	"github.com/go-chi/chi"

	"github.com/herwando/mini-wallet/module/wallet/entity/interface/handler"
	"github.com/herwando/mini-wallet/module/wallet/handler/middlewares"
)

type moduleHandler struct {
	httpHandler    *handler.Handler
	authMiddleware *middlewares.Module
}

func newRoutes(mHandler moduleHandler) *chi.Mux {
	var (
		httpHandler = mHandler.httpHandler
	)

	router := chi.NewRouter()
	router.Route("/api/v1", func(v1 chi.Router) {
		v1.Post("/init", httpHandler.AccountHandler.Init)

		v1.Route("/wallet", func(wallet chi.Router) {
			wallet.Use(
				mHandler.authMiddleware.Handler,
			)

			wallet.Post("/", httpHandler.WalletHandler.Enabled)
			wallet.Patch("/", httpHandler.WalletHandler.Disable)
			wallet.Get("/", httpHandler.WalletHandler.GetWallet)
		})
	})

	return router
}
