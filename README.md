# Eth Blockchain Parser

# API Requests
- 3 Apis is ready via postman sharable link:
Please refer to email 

OR

- APIs listing:
1. URL: 127.0.0.1:8080/current-block
   Method: GET 
   Variable: -
2. URL: 127.0.0.1:8080/address/subscribe?address=0x1ab4973a48dc892cd9971ece8e01dcc7688f8f23
   Method: GET 
   Variable: address
3. URL: 127.0.0.1:8080/transactions/list?address=0x1ab4973a48dc892cd9971ece8e01dcc7688f8f23
   Method: GET
   Variable: address

# Data Variable Explanation
impl.EthTransactionCache (data in shards)
- local memory cache about all the eth transaction pull from the third party
- get related eth transactions by using blockchain block number.

impl.ProcessedBlockCache (data in shards)
- local memory cache related block number that pull from the third party.
- info regarding the block expire period.
- use to control the memory size of impl.EthTransactionCache listing cache.

impl.ProcessedBlockSeqCache
- use to delete block number in asc sequence.
- this variable will use together with variable impl.EthTransactionCache and impl.ProcessedBlockCache.

# Future Enhancement
1. Security as being server providers.
2. Memory cache controller to avoid system being memory overused. 
3. Config setting in code for flexibility system configuration.
4. Pool Worker for background processes.
5. System need to do stress test and benchmarking so that it can be survived from system heavy usage.