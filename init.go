package web

import core "github.com/procyon-projects/procyon-core"

func init() {
	/* Initialize Pools */
	initRequestContextPool()
	initApplicationContextPool()
	initWebTransactionContextPool()
	/* Patch Matcher */
	core.Register(NewSimplePathMatcher)
	/* Request Mapping Registry*/
	core.Register(NewRequestMappingRegistry)
	/* Handler Adapter */
	core.Register(NewRequestMappingHandlerAdapter)
	/* Handler Adapter Processor */
	core.Register(NewRequestMappingHandlerAdapterProcessor)
	/* Handler Mapping */
	core.Register(NewRequestHandlerMapping)
	/* Request Handler Mapping Processor */
	core.Register(NewRequestHandlerMappingProcessor)
}
