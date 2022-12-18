package main

import (
	"fmt"
	"net/http"
	"os"
	"strings"
	"sync"

	"github.com/julienschmidt/httprouter"
	"github.com/subosito/gotenv"
)

var (
	once sync.Once
)

func main() {
	once.Do(func() {
		loadEnv()
		_ = getDBConnection()
	})

	port := os.Getenv("APP_PORT")
	router := httprouter.New()
	//router.GET("/healthz", swHandler.Healthz)

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
