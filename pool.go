package web

import "sync"

var (
	requestContextPool        sync.Pool
	webTransactionContextPool sync.Pool
	applicationContextPool    sync.Pool
)

func initRequestContextPool() {
	requestContextPool = sync.Pool{
		New: newWebRequestContext,
	}
}

func initWebTransactionContextPool() {
	webTransactionContextPool = sync.Pool{
		New: newWebTransactionContext,
	}
}

func initApplicationContextPool() {
	applicationContextPool = sync.Pool{
		New: newApplicationContext,
	}
}
