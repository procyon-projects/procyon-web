package web

import (
	"fmt"
	"github.com/google/uuid"
	context "github.com/procyon-projects/procyon-context"
	"net/http"
	"runtime/debug"
)

type Router interface {
	DoGet(res HttpResponse, req HttpRequest) error
	DoPost(res HttpResponse, req HttpRequest) error
	DoPatch(res HttpResponse, req HttpRequest) error
	DoPut(res HttpResponse, req HttpRequest) error
	DoDelete(res HttpResponse, req HttpRequest) error
	DoService(res HttpResponse, req HttpRequest) error
	DoDispatch(res HttpResponse, req HttpRequest) error
	GetHandlerChain(req HttpRequest) (*HandlerChain, error)
	GetHandlerAdapter(handler interface{}) (HandlerAdapter, error)
}

const ApplicationContextAttribute = "WEB_APPLICATION_CONTEXT"

type SimpleRouter struct {
	context         WebApplicationContext
	handlerMappings []HandlerMapping
	handlerAdapters []HandlerAdapter
}

func NewSimpleRouter(context WebApplicationContext) *SimpleRouter {
	return &SimpleRouter{
		context:         context,
		handlerMappings: make([]HandlerMapping, 0),
		handlerAdapters: make([]HandlerAdapter, 0),
	}
}

func (router *SimpleRouter) DoGet(res HttpResponse, req HttpRequest) error {
	return router.DoService(res, req)
}

func (router *SimpleRouter) DoPost(res HttpResponse, req HttpRequest) error {
	return router.DoService(res, req)
}

func (router *SimpleRouter) DoPatch(res HttpResponse, req HttpRequest) error {
	return router.DoService(res, req)
}

func (router *SimpleRouter) DoPut(res HttpResponse, req HttpRequest) error {
	return router.DoService(res, req)
}

func (router *SimpleRouter) DoDelete(res HttpResponse, req HttpRequest) error {
	return router.processRequest(res, req)
}

func (router *SimpleRouter) processRequest(res HttpResponse, req HttpRequest) error {
	return router.DoService(res, req)
}

func (router *SimpleRouter) DoService(res HttpResponse, req HttpRequest) error {
	mainContext := router.context.(ConfigurableWebApplicationContext)
	var transactionContext context.Context

	// clone the logger for transaction context
	logger := mainContext.GetLogger()
	defer func() {
		if r := recover(); r != nil {
			res.responseWriter.WriteHeader(http.StatusBadRequest)
			// when you're done with the instances, put them into pool
			httpRequestPool.Put(req)
			httpResponsePool.Put(res)
			webTransactionContextPool.Put(transactionContext)
			logger.Error(transactionContext, fmt.Sprintf("%s\n%s", r, string(debug.Stack())))
		}
	}()

	contextId, err := uuid.NewUUID()
	if err != nil {
		panic(err)
	}

	transactionContext, err = prepareWebTransactionContext(contextId,
		router.context.(context.ConfigurableContext),
	)
	if err != nil {
		panic(err)
	}
	req.AddAttribute(ApplicationContextAttribute, transactionContext)

	err = router.DoDispatch(res, req)
	if err != nil {
		panic(err)
	}
	return nil
}

func (router *SimpleRouter) DoDispatch(res HttpResponse, req HttpRequest) error {
	executionChain, err := router.GetHandlerChain(req)
	if err != nil {
		return err
	}
	executionChain.applyHandleBefore(res, req)

	var adapter HandlerAdapter
	adapter, err = router.GetHandlerAdapter(executionChain.getHandler())
	if err != nil {
		return err
	}
	adapter.Handle(executionChain.getHandler(), res, req)

	executionChain.applyHandleAfter(res, req)
	return nil
}

func (router *SimpleRouter) GetHandlerChain(req HttpRequest) (*HandlerChain, error) {
	mappings := router.handlerMappings
	if len(mappings) > 0 {
		for _, mapping := range mappings {
			chain := mapping.GetHandlerChain(req)
			if chain != nil {
				return chain, nil
			}
		}
	}
	return nil, NewNoHandlerFoundError("Request handler not  found")
}

func (router *SimpleRouter) GetHandlerAdapter(handler interface{}) (HandlerAdapter, error) {
	adapters := router.handlerAdapters
	if len(adapters) > 0 {
		for _, adapter := range adapters {
			if adapter.Supports(handler) {
				return adapter, nil
			}
		}
	}
	return nil, NewRouterError("Router handler adapter not found")
}

func (router *SimpleRouter) AddHandlerMappings(mappings ...HandlerMapping) {
	router.handlerMappings = append(router.handlerMappings, mappings...)
}

func (router *SimpleRouter) AddHandlerAdapters(adapters ...HandlerAdapter) {
	router.handlerAdapters = append(router.handlerAdapters, adapters...)
}
