package impl

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/eth-blockchain-parser/utils/cache/lru"
	"github.com/eth-blockchain-parser/utils/containers/maps"
	"github.com/eth-blockchain-parser/utils/containers/slice"
	"log"
	"strconv"
	"sync"
	"time"

	"github.com/eth-blockchain-parser/models"
	"github.com/eth-blockchain-parser/utils"
)

var (
	EthTransactionCache    = maps.New[int64, []*models.EthTransaction](new(maps.Shards[int64, []*models.EthTransaction]))
	ProcessedBlockCache    = maps.New[int64, *models.ProcessedBlockInfo](new(maps.Shards[int64, *models.ProcessedBlockInfo]))
	ProcessedBlockSeqCache = slice.New[int64](new(slice.RWMutex[int64]))
)

type EthBlockchainParserImpl struct{}

func (rm *EthBlockchainParserImpl) GetCurrentBlock(ctx context.Context) int64 {

	blockHex := callEthBlockNumber()
	blockNumber, err := strconv.ParseInt(blockHex[2:], 16, 64) // Convert hex to decimal
	if err != nil {
		log.Default().Printf("strconv.ParseInt err: %v, blockHex: %v", err, blockHex)
		return 0
	}

	return blockNumber
}

func callEthBlockNumber() string {
	responseBody, err := utils.SendEthRPCRequest("eth_blockNumber", []interface{}{})
	if err != nil {
		log.Default().Printf("utils.SendEthRPCRequest err: eth_blockNumber %v", err)
		return ""
	}

	var result map[string]interface{}
	if err := json.Unmarshal(responseBody, &result); err != nil {
		log.Default().Printf("json.Unmarshal err:%v, result:%v", err, result)
		return ""
	}
	log.Default().Printf("json.Unmarshal successfully result:%v", result)
	return result["result"].(string)
}

func (rm *EthBlockchainParserImpl) Subscribe(ctx context.Context, address string) bool {
	log.Default().Printf("Start Subscribe Address:%v", address)
	currentBlockNumber := rm.GetCurrentBlock(ctx)
	if models.SubscribedAddressesCache == nil {
		models.SubscribedAddressesCache = lru.New[string, *models.EthAddress](models.EthAddressModel.CacheSize)
	}
	if _, exists := models.SubscribedAddressesCache.Get(address); !exists {
		models.SubscribedAddressesCache.Put(address, &models.EthAddress{
			Address:            address,
			StartBlockNumber:   currentBlockNumber,
			CurrentBlockNumber: currentBlockNumber,
			CreatedAt:          time.Now(),
		})
	}
	log.Default().Printf("End Subscribe Address:%v", address)
	return true
}

type Block struct {
	Number       string                   `json:"number"`
	Transactions []*models.EthTransaction `json:"transactions"`
}

type BlockResponse struct {
	Jsonrpc string `json:"jsonrpc"`
	ID      int    `json:"id"`
	Result  Block  `json:"result"`
}

func (rm *EthBlockchainParserImpl) GetTransactions(ctx context.Context, address string) ([]*models.EthTransaction, error) {
	log.Default().Printf("Start GetTransactions Address:%v:", address)
	var transactions []*models.EthTransaction
	addressCache, exists := models.SubscribedAddressesCache.Get(address)
	if !exists {
		return transactions, fmt.Errorf("address %v not exists", address)
	}

	if len(addressCache.BlockNumberList) < 1 {
		log.Default().Printf("address %v no transaction records", address)
		return transactions, nil
	}

	countNewList := 0
	newList := make([]int64, len(addressCache.BlockNumberList))
	for i := 0; i < len(addressCache.BlockNumberList); i++ {
		ethTransactionCache, ok := EthTransactionCache.Get(addressCache.BlockNumberList[i])
		if !ok {
			// start perform delete data in EthTransactionCache since no data is exists
			continue
		}
		for _, tranx := range ethTransactionCache {
			if tranx.From == address || tranx.To == address {
				transactions = append(transactions, tranx)
			}
		}
		newList[countNewList] = addressCache.BlockNumberList[i]
		countNewList++
	}
	addressCache.BlockNumberList = newList

	log.Default().Printf("End GetTransactions Address:%v:", address)
	return transactions, nil
}

func callGetBlockByNumber(blockNum int64) *BlockResponse {
	var blockHex string
	if blockNum > 0 {
		blockHex = fmt.Sprintf("0x%x", blockNum)
	} else {
		blockHex = callEthBlockNumber()
		if blockHex == "" {
			log.Default().Printf("callEthBlockNumber err: missing blockHex")
			return nil
		}
	}
	params := []interface{}{blockHex, true}
	response, err := utils.SendEthRPCRequest("eth_getBlockByNumber", params)
	if err != nil {
		log.Default().Printf("utils.SendEthRPCRequest eth_getBlockByNumber err:  %v", err)
		return nil
	}

	var result BlockResponse
	if err := json.Unmarshal(response, &result); err != nil {
		log.Default().Printf("json.Unmarshal err: %v data: %v ", err, response)
		return nil
	}

	return &result
}

func SetupTestData(ctx context.Context) {
	blockHex := "0x13C80D9"
	params := []interface{}{blockHex, true}
	response, err := utils.SendEthRPCRequest("eth_getBlockByNumber", params)
	if err != nil {
		log.Default().Printf("utils.SendEthRPCRequest eth_getBlockByNumber err:  %v", err)
		return
	}

	var result BlockResponse
	if err := json.Unmarshal(response, &result); err != nil {
		log.Default().Printf("json.Unmarshal err: %v data: %v ", err, response)
		return
	}

	wg := sync.WaitGroup{}

	wg.Add(len(result.Result.Transactions))
	for _, tx := range result.Result.Transactions {
		blockNum, _ := strconv.ParseInt(result.Result.Number[2:], 16, 64) // Convert hex to decimal
		go func(blockNum int64, tx *models.EthTransaction) {
			defer wg.Done()
			setupTestAddr(blockNum, tx.From)
			setupTestAddr(blockNum, tx.To)
		}(blockNum, tx)
	}
	wg.Wait()
}

func setupTestAddr(blockNumber int64, address string) {
	log.Default().Printf("Start Subscribe Address:%v", address)
	if models.SubscribedAddressesCache == nil {
		models.SubscribedAddressesCache = lru.New[string, *models.EthAddress](models.EthAddressModel.CacheSize)
	}
	if _, exists := models.SubscribedAddressesCache.Get(address); !exists {
		models.SubscribedAddressesCache.Put(address, &models.EthAddress{
			Address:            address,
			StartBlockNumber:   blockNumber,
			CurrentBlockNumber: blockNumber,
			CreatedAt:          time.Now(),
		})
	}
	log.Default().Printf("End Subscribe Address:%v", address)
}

func RunTaskBg() {
	go RunScheduleUpdEthTranxBg(5 * time.Second)
	go RunScheduleUpdEthTranxBg(10 * time.Second)
	go RunScheduleRemoveExpireBlockTransaction(3 * time.Minute)
}

func RunScheduleUpdEthTranxBg(period time.Duration) {
	for range time.Tick(period) {
		RunTaskUpdEthTranxBg()
	}
}

func RunScheduleRemoveExpireBlockTransaction(period time.Duration) {
	for range time.Tick(period) {
		RemoveExpireBlockTransaction()
	}
}

func RemoveExpireBlockTransaction() error {
	if ProcessedBlockSeqCache.Len() < 1 {
		errMsg := "processed block cache is empty"
		log.Default().Println(errMsg)
		return fmt.Errorf(errMsg)
	}
	deleteCount := ProcessedBlockSeqCache.Len()
	if deleteCount > 9 {
		deleteCount = 9
	}
	dtNow := time.Now()
	for i := deleteCount; i >= 0; i-- {
		blockCount, ok := ProcessedBlockSeqCache.Index(i)
		if !ok {
			errMsg := "no processed block seq in cache. something is went wrong in cache"
			log.Default().Println(errMsg)
			return fmt.Errorf(errMsg)
		}
		if blockCount < 1 {
			errMsg := "invalid process block number in processed block seq cache"
			log.Default().Println(errMsg)
			return fmt.Errorf(errMsg)
		}

		processedBlockCacheInfo, ok := ProcessedBlockCache.Get(blockCount)
		if !ok {
			errMsg := "no processed block in cache. something is went wrong in cache"
			log.Default().Println(errMsg)
			return fmt.Errorf(errMsg)
		}
		if processedBlockCacheInfo == nil {
			errMsg := "no processed block is nil in cache. something is went wrong in cache"
			log.Default().Println(errMsg)
			return fmt.Errorf(errMsg)
		}
		if dtNow.After(processedBlockCacheInfo.ExpiredAt) {
			EthTransactionCache.Delete(blockCount)
			ProcessedBlockCache.Delete(blockCount)
			ProcessedBlockSeqCache.Delete(i)
			log.Default().Printf("block number: %v is deleted successfully", blockCount)
		} else {
			log.Default().Printf("skip delete block number: %v", blockCount)
		}
	}
	return nil
}

func RunTaskUpdEthTranxBg() error {
	ctx := context.Background()
	if models.SubscribedAddressesCache == nil {
		errMsg := "no subscribed address in cache"
		log.Default().Println(errMsg)
		return fmt.Errorf(errMsg)
	}
	result := callGetBlockByNumber(0)
	dateTime, newBlockStatus := processValidateNewBlock(result.Result.Number)
	blockNumber, err := strconv.ParseInt(result.Result.Number[2:], 16, 64) // Convert hex to decimal
	if err != nil {
		errMsg := fmt.Sprintf("block number: %v (%v) err: %v", result.Result.Number, blockNumber, err)
		log.Default().Println(errMsg)
		return fmt.Errorf(errMsg)
	}
	if !newBlockStatus {
		errMsg := fmt.Sprintf("old block: %v (%v) skip processed @ %v", result.Result.Number, blockNumber, dateTime.Format(time.DateTime))
		log.Default().Println(errMsg)
		return fmt.Errorf(errMsg)
	}
	log.Default().Printf("processing new block: %v (%v) @ %v", result.Result.Number, blockNumber, dateTime.Format(time.DateTime))
	UpdEthTranxBg(ctx, &result.Result)
	return nil
}

func UpdEthTranxBg(ctx context.Context, result *Block) {
	tmpTable := make(map[string]bool)
	var tranx []*models.EthTransaction
	blockNumber, _ := strconv.ParseInt(result.Number[2:], 16, 64) // Convert hex to decimal
	models.SubscribedAddressesCache.Range(func(address string, ethAddressInfo *models.EthAddress) error {
		models.SetLatestBlockNumber(ethAddressInfo, blockNumber)
		models.SetBlockNumberList(ethAddressInfo, blockNumber)
		for _, tx := range result.Transactions {
			if tx.From == address || tx.To == address {
				if _, ok := tmpTable[tx.Hash]; !ok {
					tmpTable[tx.Hash] = true
					tranx = append(tranx, tx)
					log.Default().Printf("address: %v in block: %v (%v) new tranx hash: %v is processed @ %v", address, result.Number, blockNumber, tx.Hash, time.Now().Format(time.DateTime))
				}
			}
			if len(tmpTable) == len(result.Transactions) {
				break
			}
		}
		return nil
	})

	//cache := ttl.New[int64, []*models.EthTransaction](models.TTLPeriod)
	//expiredAt := time.Now().Add(models.TTLPeriod).Format(time.DateTime)
	//cache.Put(blockNumber, tranx)
	EthTransactionCache.Set(blockNumber, tranx)
	log.Default().Printf("new block: %v (%v) is processed successfully", result.Number, blockNumber)
	return
}

func processValidateNewBlock(blockHex string) (time.Time, bool) {
	blockNumber, _ := strconv.ParseInt(blockHex[2:], 16, 64) // Convert hex to decimal
	processedBlockInfo, exists := ProcessedBlockCache.Get(blockNumber)
	var dateTime time.Time
	if !exists {
		dateTime = time.Now()
		ProcessedBlockSeqCache.Push(blockNumber)
		ProcessedBlockCache.Set(blockNumber, &models.ProcessedBlockInfo{
			BlockHex:    blockHex,
			BlockNumber: blockNumber,
			CreatedAt:   dateTime,
			ExpiredAt:   dateTime.Add(models.BlockTransactionModel.TTLPeriod),
		})
		return dateTime, true
	}
	dateTime = processedBlockInfo.CreatedAt
	return dateTime, false
}
