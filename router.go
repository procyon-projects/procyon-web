package web

import (
	"github.com/codnect/goo"
	context "github.com/procyon-projects/procyon-context"
	"github.com/valyala/fasthttp"
	"sync"
)

type Router interface {
	Route(requestCtx *fasthttp.RequestCtx)
}

type ProcyonRouter struct {
	ctx                context.ConfigurableApplicationContext
	handlerMapping     HandlerMapping
	requestContextPool *sync.Pool
}

func NewProcyonRouterForBenchmark(context context.ConfigurableApplicationContext, handlerRegistry SimpleHandlerRegistry) *ProcyonRouter {
	router := &ProcyonRouter{
		ctx: context,
	}
	router.handlerMapping = NewRequestHandlerMapping(NewRequestMappingRegistry())
	registryMap := handlerRegistry.getRegistryMap()
	for _, handlers := range registryMap {
		for _, handler := range handlers {
			router.handlerMapping.RegisterHandlerMethod(handler.Path, handler.Methods[0], handler.HandlerFunc)
		}
	}
	router.requestContextPool = &sync.Pool{
		New: newWebRequestContext,
	}
	return router
}

func NewProcyonRouter(context context.ConfigurableApplicationContext) *ProcyonRouter {
	router := &ProcyonRouter{
		ctx: context,
	}
	router.registerHandlerAdapter()
	return router
}

func (router *ProcyonRouter) registerHandlerAdapter() {
	handlerAdapter := router.ctx.GetSharedPeaType(goo.GetType((*HandlerMapping)(nil)))
	router.handlerMapping = handlerAdapter.(HandlerMapping)
}
func (router *ProcyonRouter) Route(requestCtx *fasthttp.RequestCtx) {
	requestContext := router.requestContextPool.Get().(*WebRequestContext)
	requestContext.fastHttpRequestContext = requestCtx

	/* if it's needed, reset the values in context */
	if requestContext.needReset {
		requestContext.reset()
	} else {
		requestContext.needReset = true
	}

	// prepare the context
	requestContext.prepare()

	// get handler chain and call all handlers
	router.handlerMapping.GetHandlerChain(requestContext)
	requestContext.Next()

	// after it's finished, put the request context to pool back
	router.requestContextPool.Put(requestContext)
}
