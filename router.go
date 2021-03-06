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
	ctx                 context.ConfigurableApplicationContext
	handlerMapping      HandlerMapping
	requestContextPool  *sync.Pool
	generateContextId   bool
	recoveryActive      bool
	errorHandlerManager *errorHandlerManager
	validator           Validator
	requestBinder       RequestBinder
	responseBodyWriter  ResponseBodyWriter
}

func newProcyonRouterForBenchmark(context context.ConfigurableApplicationContext, handlerRegistry SimpleHandlerRegistry) *ProcyonRouter {
	router := &ProcyonRouter{
		ctx: context,
	}
	router.requestContextPool = &sync.Pool{
		New: router.newWebRequestContext,
	}

	router.handlerMapping = NewRequestHandlerMapping(NewRequestMappingRegistry(), nil)
	registryMap := handlerRegistry.getRegistryMap()
	for _, handlers := range registryMap {
		for _, handler := range handlers {
			router.handlerMapping.RegisterHandlerMethod(handler.Path, handler.Method, handler.HandlerFunc, nil)
		}
	}
	return router
}

func NewProcyonRouter(context context.ConfigurableApplicationContext) *ProcyonRouter {
	router := &ProcyonRouter{
		ctx:                context,
		generateContextId:  true,
		recoveryActive:     true,
		validator:          newDefaultValidator(),
		requestBinder:      newDefaultRequestBinder(),
		responseBodyWriter: newDefaultResponseBodyWriter(),
	}
	router.requestContextPool = &sync.Pool{
		New: router.newWebRequestContext,
	}
	router.configure()
	return router
}

func (router *ProcyonRouter) newWebRequestContext() interface{} {
	requestContext := &WebRequestContext{
		router:       router,
		handlerIndex: 0,
		valueMap:     make(map[string]interface{}),
	}
	return requestContext
}

func (router *ProcyonRouter) configure() {
	peaFactory := router.ctx.GetPeaFactory()

	handlerAdapter := peaFactory.GetSharedPeaType(goo.GetType((*HandlerMapping)(nil)))
	router.handlerMapping = handlerAdapter.(HandlerMapping)

	// custom logger
	router.errorHandlerManager = newErrorHandlerManager(router.ctx.GetLogger())
	errorHandler, _ := peaFactory.GetPeaByType(goo.GetType((*ErrorHandler)(nil)))
	if errorHandler != nil {
		router.errorHandlerManager.customErrorHandler = errorHandler.(ErrorHandler)
	}

	// custom validator
	customValidator, _ := peaFactory.GetPeaByType(goo.GetType((*Validator)(nil)))
	if customValidator != nil {
		router.validator = customValidator.(Validator)
	}

	// custom request binder
	customRequestBinder, _ := peaFactory.GetPeaByType(goo.GetType((*RequestBinder)(nil)))
	if customRequestBinder != nil {
		router.requestBinder = customRequestBinder.(RequestBinder)
	}

	// custom response body writer
	customResponseBodyWriter, _ := peaFactory.GetPeaByType(goo.GetType((*ResponseBodyWriter)(nil)))
	if customResponseBodyWriter != nil {
		router.responseBodyWriter = customResponseBodyWriter.(ResponseBodyWriter)
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
		router.errorHandlerManager.HandleError(HttpErrorNotFound, requestContext)

		requestContext.reset()
		router.requestContextPool.Put(requestContext)
		return
	}

	requestContext.invoke()

	requestContext.reset()
	router.requestContextPool.Put(requestContext)
}
