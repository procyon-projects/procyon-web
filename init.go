package web

import core "github.com/procyon-projects/procyon-core"

func init() {
	/* Initialize Pools */
	initHttpRequestPool()
	initHttpResponsePool()
	initWebTransactionContextPool()
	/* Request Handler Mapping */
	core.Register(NewRequestHandlerMapping)
}
