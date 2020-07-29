package web

import "sync"

var (
	httpRequestPool        sync.Pool
	httpResponsePool       sync.Pool
	transactionContextPool sync.Pool
)

func initHttpRequestPool() {
	httpRequestPool = sync.Pool{
		New: newHttpRequest,
	}
}

func initHttpResponsePool() {
	httpResponsePool = sync.Pool{
		New: newHttpResponse,
	}
}

func initTransactionContextPool() {
	transactionContextPool = sync.Pool{
		New: newTransactionContext,
	}
}
