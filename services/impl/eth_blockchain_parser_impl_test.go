package impl

import (
	"context"
	"fmt"
	"testing"
)

func TestGetCurrentBlock(t *testing.T) {

	parser := new(EthBlockchainParserImpl)
	block := parser.GetCurrentBlock(context.TODO())
	fmt.Println(block)
}

func TestSubscribe(t *testing.T) {

	parser := new(EthBlockchainParserImpl)
	fmt.Println(parser.Subscribe(context.TODO(), "0xa264e3f3ef78f5bb93476fb53a9617558e303142"))
}

func TestGetTransactions(t *testing.T) {

	parser := new(EthBlockchainParserImpl)
	fmt.Println(parser.GetTransactions(context.TODO(), "0xa264e3f3ef78f5bb93476fb53a9617558e303142"))
}

func TestUpdEthTranxBg(t *testing.T) {
	ctx := context.Background()
	SetupTestData(ctx)
	result := callGetBlockByNumber(0)
	UpdEthTranxBg(ctx, &result.Result)
}

func TestRunTaskBg(t *testing.T) {
	ctx := context.Background()
	SetupTestData(ctx)
	RunTaskBg()
}
