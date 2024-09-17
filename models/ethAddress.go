package models

import (
	"github.com/eth-blockchain-parser/utils/array"
	"github.com/eth-blockchain-parser/utils/cache/lru"
	"time"
)

type EthAddress struct {
	Address            string
	StartBlockNumber   int64
	CurrentBlockNumber int64
	BlockNumberList    []int64 // a list of block number that exists the transaction for the related address.
	CreatedAt          time.Time
}

type ethAddressModel struct {
	CacheSize int
}

// EthAddressModel
var (
	EthAddressModel = ethAddressModel{
		CacheSize: 10000,
	}
	SubscribedAddressesCache *lru.Cache[string, *EthAddress]
)

func SetLatestBlockNumber(m *EthAddress, crtBlockNumber int64) {
	if m.CurrentBlockNumber < crtBlockNumber {
		m.CurrentBlockNumber = crtBlockNumber
	}
}

func SetBlockNumberList(m *EthAddress, crtBlockNumber int64) {
	if len(m.BlockNumberList) < 1 || m.BlockNumberList == nil {
		m.BlockNumberList = []int64{crtBlockNumber}
		return
	}
	if !array.In(m.BlockNumberList, crtBlockNumber) {
		m.BlockNumberList = append(m.BlockNumberList, crtBlockNumber)
	}
}
