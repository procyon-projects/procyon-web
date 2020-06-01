package web

type Router interface {
	DoGet(res HttpResponse, req HttpRequest) error
	DoPost(res HttpResponse, req HttpRequest) error
	DoPatch(res HttpResponse, req HttpRequest) error
	DoPut(res HttpResponse, req HttpRequest) error
	DoDelete(res HttpResponse, req HttpRequest) error
	DoService(res HttpResponse, req HttpRequest) error
	DoDispatch(res HttpResponse, req HttpRequest) error
	GetHandlerChain(req HttpRequest) (*HandlerChain, error)
	GetHandlerAdapter(req HttpRequest) (*HandlerAdapter, error)
}

const ApplicationContextAttribute = "DEFAULT_HANDLER_CONTEXT"

type SimpleRouter struct {
	context         ApplicationContext
	handlerMappings []HandlerMapping
	handlerAdapters []HandlerAdapter
}

func NewSimpleRouter() *SimpleRouter {
	return &SimpleRouter{
		context:         nil,
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
	/* logging etc... */
	req.AddAttribute(ApplicationContextAttribute, router.context)
	_ = router.DoDispatch(res, req)
	return nil
}

func (router *SimpleRouter) DoDispatch(res HttpResponse, req HttpRequest) error {
	executionChain, err := router.GetHandlerExecutionChain(req)
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

func (router *SimpleRouter) GetHandlerExecutionChain(req HttpRequest) (*HandlerChain, error) {
	if len(router.handlerMappings) > 0 {
		for _, handlerMapping := range router.handlerMappings {
			chain := handlerMapping.GetHandlerChain(req)
			if chain != nil {
				return chain, nil
			}
		}
	}
	return nil, NewNoHandlerFoundError("Request handler not  found")
}

func (router *SimpleRouter) GetHandlerAdapter(handler interface{}) (HandlerAdapter, error) {
	if len(router.handlerAdapters) > 0 {
		for _, handlerAdapter := range router.handlerAdapters {
			if handlerAdapter.Supports(handler) {
				return handlerAdapter, nil
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
