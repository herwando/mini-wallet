package main

import (
	"github.com/go-chi/chi"

	"github.com/herwando/mini-wallet/module/wallet/entity/interface/handler"
)

type moduleHandler struct {
	httpHandler *handler.Handler
}

func newRoutes(mHandler moduleHandler) *chi.Mux {
	var (
		walletHandler = mHandler.walletHandler
	)

	router := chi.NewRouter()
	router.Get("/ping", walletHandler.Ping)

	return router
}
