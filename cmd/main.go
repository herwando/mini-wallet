package main

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	httpHandler "github.com/herwando/mini-wallet/module/wallet/entity/interface/handler"
	"github.com/herwando/mini-wallet/module/wallet/handler"
	"github.com/herwando/mini-wallet/module/wallet/handler/middlewares"
	"github.com/herwando/mini-wallet/module/wallet/repository"
	"github.com/herwando/mini-wallet/module/wallet/usecase"
	"github.com/subosito/gotenv"
)

func main() {
	loadEnv()
	db := getDBConnection()

	walletRepo := repository.NewWalletRepository(db)
	walletUsecase := usecase.NewWalletUsecase(walletRepo)
	walletHandler := handler.NewWalletHandler(walletUsecase)
	accountRepo := repository.NewAccountRepository(db)
	accountUsecase := usecase.NewAccountUsecase(accountRepo)
	accountHandler := handler.NewAccountHandler(accountUsecase)
	depositRepo := repository.NewDepositRepository(db)
	depositUsecase := usecase.NewDepositUsecase(depositRepo, walletRepo)
	depositHandler := handler.NewDepositHandler(depositUsecase)
	withdrawalRepo := repository.NewWithdrawalRepository(db)
	withdrawalUsecase := usecase.NewWithdrawalUsecase(withdrawalRepo, walletRepo)
	withdrawalHandler := handler.NewWithdrawalHandler(withdrawalUsecase)

	port := os.Getenv("APP_PORT")
	handler := &httpHandler.Handler{
		WalletHandler:     walletHandler,
		AccountHandler:    accountHandler,
		DepositHandler:    depositHandler,
		WithdrawalHandler: withdrawalHandler,
	}
	authHandler := &middlewares.Module{}

	router := newRoutes(moduleHandler{
		httpHandler:    handler,
		authMiddleware: authHandler,
	})

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
