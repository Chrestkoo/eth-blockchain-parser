package models

type EthTransaction struct {
	BlockHash        string `json:"blockHash" bson:"blockHash"`
	BlockNumber      string `json:"blockNumber" bson:"blockNumber"`
	ChainID          string `json:"chainId" bson:"chainId"`
	From             string `json:"from" bson:"from"`
	Gas              string `json:"gas" bson:"gas"`
	GasPrice         string `json:"gasPrice" bson:"gasPrice"`
	Hash             string `json:"hash" bson:"hash"`
	Input            string `json:"input" bson:"input"`
	Nonce            string `json:"nonce" bson:"nonce"`
	R                string `json:"r" bson:"r"`
	S                string `json:"s" bson:"s"`
	To               string `json:"to" bson:"to"`
	TransactionIndex string `json:"transactionIndex" bson:"transactionIndex"`
	Type             string `json:"type" bson:"type"`
	V                string `json:"v" bson:"v"`
	Value            string `json:"value" bson:"value"`
}

type ethTransactionModel struct {
	CacheSize int
}

// EthTransactionModel
var (
	EthTransactionModel = ethTransactionModel{
		CacheSize: 5000,
	}
)
