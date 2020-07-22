package web

import "github.com/google/uuid"

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

const ApplicationContextAttribute = "SIMPLE_ROUTER_WEB_CONTEXT"

type SimpleRouter struct {
	context         ApplicationContext
	handlerMappings []HandlerMapping
	handlerAdapters []HandlerAdapter
}

func NewSimpleRouter(context ApplicationContext) *SimpleRouter {
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
	mainContext := router.context.(ConfigurableApplicationContext)

	// clone the logger for transaction context
	logger := mainContext.GetLogger()
	contextId, err := uuid.NewUUID()
	if err != nil {
		logger.Panic(err)
		return nil
	}
	transactionLogger := logger.Clone(contextId)

	defer func() {
		if r := recover(); r != nil {
			transactionLogger.Panic(r)
		}
	}()

	var txContext *TransactionContext
	txContext, err = prepareTransactionContext(contextId, router.context.(ConfigurableApplicationContext), logger)
	if err != nil {
		transactionLogger.Panic(err)
		return nil
	}
	req.AddAttribute(ApplicationContextAttribute, txContext)

	err = router.DoDispatch(res, req)
	if err != nil {
		transactionLogger.Panic(err)
		return nil
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
