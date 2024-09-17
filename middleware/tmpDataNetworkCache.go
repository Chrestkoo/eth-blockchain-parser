package middleware

// TODO
/*
	GetTmpDataCache

	This function is to reduce the load from the cache data (impl.EthTransactionCache).
	This use ttl for the temporary cache.
	The data shall delete after some time when data is expired.
	Once the data is expired, the request get the data from the real cache data (impl.EthTransactionCache).
*/
func GetTmpDataCache() {

}
