package utils

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
)

const ethNodeURL = "https://cloudflare-eth.com"

// JSON-RPC request struct
type JsonRPCRequest struct {
	Jsonrpc string        `json:"jsonrpc"`
	Method  string        `json:"method"`
	Params  []interface{} `json:"params"`
	ID      int           `json:"id"`
}

// Function to send eth JSON-RPC request
func SendEthRPCRequest(method string, params []interface{}) ([]byte, error) {
	reqBody, err := json.Marshal(JsonRPCRequest{
		Jsonrpc: "2.0",
		Method:  method,
		Params:  params,
		ID:      1,
	})
	if err != nil {
		return nil, err
	}

	resp, err := http.Post(ethNodeURL, "application/json", bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return io.ReadAll(resp.Body)
}
