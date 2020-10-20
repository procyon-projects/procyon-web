package web

import core "github.com/procyon-projects/procyon-core"

func init() {
	/* Initialize Pools */
	initHttpRequestPool()
	initHttpResponsePool()
	initApplicationContextPool()
	initWebTransactionContextPool()
	/* Request Handler Mapping Processor */
	core.Register(NewRequestHandlerMappingProcessor)
}
