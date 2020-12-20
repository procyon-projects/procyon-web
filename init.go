package web

import core "github.com/procyon-projects/procyon-core"

func init() {
	/* Request Mapping Registry*/
	core.Register(NewRequestMappingRegistry)
	/* Handler Mapping */
	core.Register(NewRequestHandlerMapping)
	/* Request Handler Mapping Processor */
	core.Register(NewRequestHandlerMappingProcessor)
	/* Handler Interceptor Registry & Processor */
	core.Register(NewSimpleHandlerInterceptorRegistry)
	core.Register(NewHandlerInterceptorProcessor)
}
