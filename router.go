package web

import (
	"github.com/procyon-projects/goo"
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
	recoveryManager    *recoveryManager
}

func newProcyonRouterForBenchmark(context context.ConfigurableApplicationContext, handlerRegistry SimpleHandlerRegistry) *ProcyonRouter {
	router := &ProcyonRouter{
		ctx: context,
		requestContextPool: &sync.Pool{
			New: newWebRequestContext,
		},
	}
	router.handlerMapping = NewRequestHandlerMapping(NewRequestMappingRegistry(), nil)
	registryMap := handlerRegistry.getRegistryMap()
	for _, handlers := range registryMap {
		for _, handler := range handlers {
			router.handlerMapping.RegisterHandlerMethod(handler.Path, handler.Method, handler.HandlerFunc)
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
	router.configure()
	return router
}

func (router *ProcyonRouter) configure() {
	peaFactory := router.ctx.GetPeaFactory()

	handlerAdapter := peaFactory.GetSharedPeaType(goo.GetType((*HandlerMapping)(nil)))
	router.handlerMapping = handlerAdapter.(HandlerMapping)

	router.recoveryManager = newRecoveryManager(router.ctx.GetLogger())
	errorHandler, _ := peaFactory.GetPeaByType(goo.GetType((*ErrorHandler)(nil)))
	if errorHandler != nil {
		router.recoveryManager.customErrorHandler = errorHandler.(ErrorHandler)
	}
}

func (router *ProcyonRouter) Route(requestCtx *fasthttp.RequestCtx) {
	requestContext := router.requestContextPool.Get().(*WebRequestContext)
	requestContext.fastHttpRequestContext = requestCtx
	// prepare the context
	requestContext.prepare(router.generateContextId)

	// get handler chain and call all handlers
	router.handlerMapping.GetHandlerChain(requestContext)

	if requestContext.handlerChain == nil {
		router.ctx.GetLogger().Warning(requestContext, "Handler not found : "+string(requestCtx.Path()))
		router.recoveryManager.HandleError(HttpErrorNotFound, requestContext)

		requestContext.reset()
		router.requestContextPool.Put(requestContext)
		return
	}

	requestContext.invoke(router.recoveryActive, router.recoveryManager)
	requestContext.reset()

	// after it's finished, put the request context to pool back
	router.requestContextPool.Put(requestContext)
}
