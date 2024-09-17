package main

import (
	"fmt"
	"github.com/eth-blockchain-parser/middleware"
	"github.com/eth-blockchain-parser/services/impl"
	"log"
	"net/http"

	"github.com/eth-blockchain-parser/controller"
)

func main() {
	//impl.SetupTestData(context.Background())
	go impl.RunTaskBg()
	initHttp()
}

func initHttp() {
	mux := http.NewServeMux()
	portNum := 8080
	log.Default().Printf("start initializing http...")

	mux.HandleFunc("/current-block", controller.GetCurrentBlockHandler)
	mux.HandleFunc("/address/subscribe", controller.SubscribeHandler)
	mux.HandleFunc("/transactions/list", controller.GetTransactionsHandler)

	// Start the cleanup goroutine
	go middleware.CleanupIPs()

	log.Default().Printf("listening localhost:%v", portNum)
	err := http.ListenAndServe(fmt.Sprintf(":%v", portNum), middleware.RateLimiterMiddleware(mux))
	if err != nil {
		log.Default().Printf("http.ListenAndServe err:%v at portNum: portNum", err, portNum)
	}
}
