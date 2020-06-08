package web

type HandlerAdapter interface {
	Supports(handler interface{}) bool
	Handle(handler interface{}, res HttpResponse, req HttpRequest) interface{}
}

type RequestMappingHandlerAdapter struct {
	parameterResolvers  *HandlerMethodParameterResolvers
	returnValueHandlers *HandlerMethodReturnValueHandlers
}

type RequestMappingHandlerAdapterOption func(adapter *RequestMappingHandlerAdapter)

func NewRequestMappingHandlerAdapter(options ...RequestMappingHandlerAdapterOption) *RequestMappingHandlerAdapter {
	adapter := &RequestMappingHandlerAdapter{
		parameterResolvers:  getDefaultMethodParameterResolvers(),
		returnValueHandlers: getDefaultReturnValueHandlers(),
	}
	for _, option := range options {
		option(adapter)
	}
	return adapter
}

func WithCustomParameterResolvers(resolvers ...HandlerMethodParameterResolver) RequestMappingHandlerAdapterOption {
	return func(adapter *RequestMappingHandlerAdapter) {
		adapter.parameterResolvers.AddMethodParameterResolver(resolvers...)
	}
}

func WithCustomReturnValueHandlers(handlers ...HandlerMethodReturnValueHandler) RequestMappingHandlerAdapterOption {
	return func(adapter *RequestMappingHandlerAdapter) {
		adapter.returnValueHandlers.AddMethodReturnValueHandler(handlers...)
	}
}

func getDefaultMethodParameterResolvers() *HandlerMethodParameterResolvers {
	resolvers := NewHandlerMethodParameterResolvers()
	resolvers.AddMethodParameterResolver(NewRequestMethodParameterResolver())
	return resolvers
}

func getDefaultReturnValueHandlers() *HandlerMethodReturnValueHandlers {
	handlers := NewHandlerMethodReturnValueHandlers()
	handlers.AddMethodReturnValueHandler(
		NewResponseEntityReturnValueHandler(),
		NewErrorReturnValueHandler(),
	)
	return handlers
}

func (adapter RequestMappingHandlerAdapter) Supports(handler interface{}) bool {
	if _, ok := handler.(HandlerMethod); ok {
		return true
	}
	return false
}

func (adapter RequestMappingHandlerAdapter) Handle(handler interface{}, res HttpResponse, req HttpRequest) interface{} {
	return adapter.invokeHandler(handler.(HandlerMethod), res, req)
}

func (adapter RequestMappingHandlerAdapter) invokeHandler(handler HandlerMethod, res HttpResponse, req HttpRequest) interface{} {
	return nil
}
