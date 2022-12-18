package main

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/herwando/mini-wallet/module/wallet/handler"
	"github.com/herwando/mini-wallet/module/wallet/repository"
	"github.com/herwando/mini-wallet/module/wallet/usecase"
	"github.com/julienschmidt/httprouter"
	"github.com/subosito/gotenv"
)

func main() {
	loadEnv()
	db := getDBConnection()

	walletRepo := repository.NewWalletRepository(db)
	walletUsecase := usecase.NewWalletUsecase(walletRepo)
	walletHandler := handler.NewWalletHandler(walletUsecase)

	port := os.Getenv("APP_PORT")
	router := httprouter.New()
	router.GET("/ping", walletHandler.Ping)

	fmt.Println("mini-wallet is now running and ready to listen at port", port)
	err := http.ListenAndServe(":"+port, router)
	fmt.Println("mini-wallet error:", err)
}

func loadEnv() {
	environment, ok := os.LookupEnv("ENVIRONMENT")
	if !ok {
		environment = "DEVELOPMENT"
	}
	_ = os.Setenv("ENVIRONMENT", strings.ToUpper(strings.TrimSpace(environment)))
	_ = gotenv.Load()
}
