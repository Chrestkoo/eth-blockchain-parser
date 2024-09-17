package impl

import (
	"context"

	"github.com/eth-blockchain-parser/models"
)

type Parser interface {
	// last parsed block
	GetCurrentBlock(ctx context.Context) int

	// add address to observer
	Subscribe(ctx context.Context, address string) bool

	// list of inbound or outbound transactions for an address
	GetTransactions(ctx context.Context, address string) []*models.EthTransaction
}
