package web

import core "github.com/procyon-projects/procyon-core"

func init() {
	/* Initialize Pools */
	initHttpRequestPool()
	initHttpResponsePool()
	initApplicationContextPool()
	initWebTransactionContextPool()
	/* Patch Matcher */
	core.Register(NewSimplePathMatcher)
	/* Request Mapping Registry*/
	core.Register(NewRequestMappingRegistry)
	/* Handler Mapping */
	core.Register(NewRequestHandlerMapping)
	/* Request Handler Mapping Processor */
	core.Register(NewRequestHandlerMappingProcessor)
}
