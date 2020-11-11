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
	DoDispatch(requestContext RequestContext, res http.ResponseWriter, req *http.Request)
	GetRequestHandler(requestContext RequestContext, req *http.Request) (*HandlerMethod, error)
}

type ProcyonRouter struct {
	context               ConfigurableWebApplicationContext
	interceptors          []HandlerInterceptor
	handlerMappings       []HandlerMapping
	requestMappingAdapter HandlerAdapter
}

func newProcyonRouterForBenchmark(context ConfigurableWebApplicationContext, handlerRegistry SimpleHandlerRegistry) *ProcyonRouter {
	router := &ProcyonRouter{
		context:         context,
		interceptors:    make([]HandlerInterceptor, 0),
		handlerMappings: make([]HandlerMapping, 0),
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
	router.SetHandlerAdapter(NewRequestMappingHandlerAdapter(core.NewDefaultTypeConverterService()))
	router.requestMappingAdapter = NewRequestMappingHandlerAdapter(core.NewDefaultTypeConverterService())
	return router
}

func NewProcyonRouter(context ConfigurableWebApplicationContext) *ProcyonRouter {
	router := &ProcyonRouter{
		context:         context,
		handlerMappings: make([]HandlerMapping, 0),
	}
	router.configureRouter()
	return router
}

func (router *ProcyonRouter) configureRouter() {
	router.registerHandlerMappings()
	router.registerHandlerAdapter()
}

func (router *ProcyonRouter) registerHandlerMappings() {
	handlerMappings := router.context.GetSharedPeasByType(goo.GetType((*HandlerMapping)(nil)))
	for _, handlerMapping := range handlerMappings {
		router.AddHandlerMappings(handlerMapping.(HandlerMapping))
	}
}

func (router *ProcyonRouter) registerHandlerAdapter() {
	handlerAdapter := router.context.GetSharedPeaType(goo.GetType((*HandlerAdapter)(nil)))
	router.SetHandlerAdapter(handlerAdapter.(HandlerAdapter))
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
	defer func() {
		if r := recover(); r != nil {
			res.WriteHeader(http.StatusBadRequest)
			requestContextPool.Put(requestContext)
		} else {
			requestContextPool.Put(requestContext)
		}
	}()
	requestContext.reset()
	requestContext.request = req
	requestContext.hasContextId = true
	router.DoDispatch(requestContext, res, req)
}

func (router *ProcyonRouter) DoDispatch(requestContext RequestContext, res http.ResponseWriter, req *http.Request) {
	handler, _ := router.GetRequestHandler(requestContext, req)

	for _, interceptor := range router.interceptors {
		interceptor.HandleBefore(handler, requestContext, res, req)
	}

	router.requestMappingAdapter.Handle(handler, requestContext, res, req)

	for _, interceptor := range router.interceptors {
		interceptor.HandleAfter(handler, requestContext, res, req)
	}

}

func (router *ProcyonRouter) GetRequestHandler(requestContext RequestContext, req *http.Request) (*HandlerMethod, error) {

	handler := router.handlerMappings[0].GetHandler(requestContext, req)
	if handler != nil {
		return handler, nil
	}

	return nil, NewNoHandlerFoundError("Request handler not  found")
}

func (router *ProcyonRouter) AddHandlerMappings(mappings ...HandlerMapping) {
	router.handlerMappings = append(router.handlerMappings, mappings...)
}

func (router *ProcyonRouter) SetHandlerAdapter(adapter HandlerAdapter) {
	router.requestMappingAdapter = adapter
}
