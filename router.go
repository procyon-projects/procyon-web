package web

import (
	"github.com/codnect/goo"
	core "github.com/procyon-projects/procyon-core"
	"net/http"
)

type handlerMetadata struct {
	HandlerName string
	Mapping     interface{}
	Fun         RequestHandlerFunc
}

type Router interface {
	DoGet(res http.ResponseWriter, req *http.Request)
	DoPost(res http.ResponseWriter, req *http.Request)
	DoPatch(res http.ResponseWriter, req *http.Request)
	DoPut(res http.ResponseWriter, req *http.Request)
	DoDelete(res http.ResponseWriter, req *http.Request)
	DoService(res http.ResponseWriter, req *http.Request)
	DoDispatch(requestContext RequestContext, res http.ResponseWriter, req *http.Request) error
	GetRequestHandler(requestContext RequestContext, req *http.Request) (HandlerMethod, error)
	GetHandlerAdapter(handler interface{}, requestContext RequestContext) (HandlerAdapter, error)
}

type ProcyonRouter struct {
	context         ConfigurableWebApplicationContext
	interceptors    []HandlerInterceptor
	handlerMappings []HandlerMapping
	handlerAdapters []HandlerAdapter
}

func newProcyonRouterForBenchmark(context ConfigurableWebApplicationContext, handlerRegistry SimpleHandlerRegistry) *ProcyonRouter {
	router := &ProcyonRouter{
		context:         context,
		interceptors:    make([]HandlerInterceptor, 0),
		handlerMappings: make([]HandlerMapping, 0),
		handlerAdapters: make([]HandlerAdapter, 0),
	}

	patchMatcher := NewSimplePathMatcher()
	requestHandlerMapping := NewRequestHandlerMapping(patchMatcher, NewRequestMappingRegistry())
	registryMap := handlerRegistry.getRegistryMap()
	for prefix, handlers := range registryMap {
		for _, handler := range handlers {
			requestMappingInfo := NewRequestMapping(newMethodRequestMatcher(handler.Methods),
				newPatternRequestMatcher(patchMatcher, prefix, handler.Paths),
			)
			requestHandlerMapping.RegisterHandlerMethod("benchmark_handler", requestMappingInfo, handler.HandlerFunc)
		}
	}
	router.AddHandlerMappings(requestHandlerMapping)
	router.AddHandlerAdapters(NewRequestMappingHandlerAdapter(core.NewDefaultTypeConverterService()))
	return router
}

func NewProcyonRouter(context ConfigurableWebApplicationContext) *ProcyonRouter {
	router := &ProcyonRouter{
		context:         context,
		handlerMappings: make([]HandlerMapping, 0),
		handlerAdapters: make([]HandlerAdapter, 0),
	}
	router.configureRouter()
	return router
}

func (router *ProcyonRouter) configureRouter() {
	router.registerHandlerMappings()
	router.registerHandlerAdapters()
}

func (router *ProcyonRouter) registerHandlerMappings() {
	handlerMappings := router.context.GetSharedPeasByType(goo.GetType((*HandlerMapping)(nil)))
	for _, handlerMapping := range handlerMappings {
		router.AddHandlerMappings(handlerMapping.(HandlerMapping))
	}
}

func (router *ProcyonRouter) registerHandlerAdapters() {
	handlerAdapters := router.context.GetSharedPeasByType(goo.GetType((*HandlerAdapter)(nil)))
	for _, handlerAdapter := range handlerAdapters {
		router.AddHandlerAdapters(handlerAdapter.(HandlerAdapter))
	}
}

func (router *ProcyonRouter) DoGet(res http.ResponseWriter, req *http.Request) {
	router.DoService(res, req)
}

func (router *ProcyonRouter) DoPost(res http.ResponseWriter, req *http.Request) {
	router.DoService(res, req)
}

func (router *ProcyonRouter) DoPatch(res http.ResponseWriter, req *http.Request) {
	router.DoService(res, req)
}

func (router *ProcyonRouter) DoPut(res http.ResponseWriter, req *http.Request) {
	router.DoService(res, req)
}

func (router *ProcyonRouter) DoDelete(res http.ResponseWriter, req *http.Request) {
	router.DoService(res, req)
}

func (router *ProcyonRouter) DoService(res http.ResponseWriter, req *http.Request) {
	requestContext := requestContextPool.Get().(*WebRequestContext)

	// clone the logger for transaction context
	//logger := router.context.GetLogger()
	defer func() {
		if r := recover(); r != nil {
			res.WriteHeader(http.StatusBadRequest)
			//applicationContextPool.Put(transactionContext.(*BaseWebApplicationContext).BaseApplicationContext)
			//webTransactionContextPool.Put(transactionContext)
			//logger.Error(transactionContext, fmt.Sprintf("%s\n%s", r, string(debug.Stack())))
			requestContext.clear()
			requestContextPool.Put(requestContext)
		} else {
			requestContext.clear()
			requestContextPool.Put(requestContext)
		}
	}()

	/*contextId, err := uuid.NewUUID()
	if err != nil {
		panic(err)
	}

	transactionContext, err = prepareWebTransactionContext(contextId,
		router.context.(context.ConfigurableContext),
	)
	if err != nil {
		panic(err)
	}*/

	err := router.DoDispatch(requestContext, res, req)
	if err != nil {
		panic(err)
	}
}

func (router *ProcyonRouter) DoDispatch(requestContext RequestContext, res http.ResponseWriter, req *http.Request) error {
	handler, err := router.GetRequestHandler(requestContext, req)
	if err != nil {
		return err
	}
	for _, interceptor := range router.interceptors {
		interceptor.HandleBefore(handler, requestContext, res, req)
	}

	var adapter HandlerAdapter
	adapter, err = router.GetHandlerAdapter(handler, requestContext)
	if err != nil {
		return err
	}
	adapter.Handle(handler, requestContext, res, req)

	for _, interceptor := range router.interceptors {
		interceptor.HandleAfter(handler, requestContext, res, req)
	}
	return nil
}

func (router *ProcyonRouter) GetRequestHandler(requestContext RequestContext, req *http.Request) (HandlerMethod, error) {
	mappings := router.handlerMappings
	if len(mappings) > 0 {
		for _, mapping := range mappings {
			handler := mapping.GetHandler(requestContext, req)
			if handler != nil {
				return handler, nil
			}
		}
	}
	return nil, NewNoHandlerFoundError("Request handler not  found")
}

func (router *ProcyonRouter) GetHandlerAdapter(handler interface{}, requestContext RequestContext) (HandlerAdapter, error) {
	adapters := router.handlerAdapters
	if len(adapters) > 0 {
		for _, adapter := range adapters {
			if adapter.Supports(handler, requestContext) {
				return adapter, nil
			}
		}
	}
	return nil, NewRouterError("Router handler adapter not found")
}

func (router *ProcyonRouter) AddHandlerMappings(mappings ...HandlerMapping) {
	router.handlerMappings = append(router.handlerMappings, mappings...)
}

func (router *ProcyonRouter) AddHandlerAdapters(adapters ...HandlerAdapter) {
	router.handlerAdapters = append(router.handlerAdapters, adapters...)
}
