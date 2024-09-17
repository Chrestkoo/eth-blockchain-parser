package models

import (
	"time"
)

type ProcessedBlockInfo struct {
	BlockHex    string    // blockchain block hex from third party
	BlockNumber int64     // blockchain block number in decimal
	CreatedAt   time.Time // block number create time
	ExpiredAt   time.Time // block number expire period. setting is 91 days.
}

type blockTransactionModel struct {
	TTLPeriod time.Duration
}

const (
	ProcessedBlockCacheSize = 5000
)

var (
	BlockTransactionModel = blockTransactionModel{
		TTLPeriod: 24 * time.Hour * 91,
	}
	//ProcessedBlockCache *lru.Cache[int64, *ProcessedBlockInfo]
)

func init() {
	//ProcessedBlockCache = lru.New[int64, *ProcessedBlockInfo](ProcessedBlockCacheSize)
}
