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
	generateContextId  bool
	recoveryActive     bool
}

func newProcyonRouterForBenchmark(context context.ConfigurableApplicationContext, handlerRegistry SimpleHandlerRegistry) *ProcyonRouter {
	router := &ProcyonRouter{
		ctx: context,
		requestContextPool: &sync.Pool{
			New: newWebRequestContext,
		},
	}
	router.handlerMapping = NewRequestHandlerMapping(NewRequestMappingRegistry())
	registryMap := handlerRegistry.getRegistryMap()
	for _, handlers := range registryMap {
		for _, handler := range handlers {
			router.handlerMapping.RegisterHandlerMethod(handler.Path, handler.Methods[0], handler.HandlerFunc)
		}
	}
	return router
}

func NewProcyonRouter(context context.ConfigurableApplicationContext) *ProcyonRouter {
	router := &ProcyonRouter{
		ctx: context,
		requestContextPool: &sync.Pool{
			New: newWebRequestContext,
		},
		generateContextId: true,
		recoveryActive:    true,
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
	// prepare the context
	requestContext.prepare(router.generateContextId)

	// get handler chain and call all handlers
	router.handlerMapping.GetHandlerChain(requestContext)
	requestContext.invoke(router.recoveryActive)
	requestContext.reset()

	// after it's finished, put the request context to pool back
	router.requestContextPool.Put(requestContext)
}
