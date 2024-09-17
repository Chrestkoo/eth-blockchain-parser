package controller

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/eth-blockchain-parser/services/impl"
)

const ApiRequestTimeout = 5 * time.Second

type APIResponse struct {
	Data interface{} `json:"data"`
}

func GetCurrentBlockHandler(w http.ResponseWriter, r *http.Request) {
	ctx, _ := context.WithTimeout(context.TODO(), ApiRequestTimeout)
	parser := new(impl.EthBlockchainParserImpl)
	block := parser.GetCurrentBlock(ctx)
	response := APIResponse{Data: block}
	json.NewEncoder(w).Encode(response)
}

func SubscribeHandler(w http.ResponseWriter, r *http.Request) {
	ctx, _ := context.WithTimeout(context.TODO(), ApiRequestTimeout)
	address := r.URL.Query().Get("address")
	if address == "" {
		http.Error(w, "Address required", http.StatusBadRequest)
		return
	}

	parser := new(impl.EthBlockchainParserImpl)
	success := parser.Subscribe(ctx, address)
	response := APIResponse{Data: success}
	json.NewEncoder(w).Encode(response)
}

func GetTransactionsHandler(w http.ResponseWriter, r *http.Request) {
	address := r.URL.Query().Get("address")
	if address == "" {
		http.Error(w, "Address required", http.StatusBadRequest)
		return
	}
	ctx, _ := context.WithTimeout(context.TODO(), ApiRequestTimeout)
	parser := new(impl.EthBlockchainParserImpl)
	transactions, _ := parser.GetTransactions(ctx, address)
	response := APIResponse{Data: transactions}
	json.NewEncoder(w).Encode(response)
}
