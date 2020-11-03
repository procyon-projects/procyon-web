package web

import "github.com/procyon-projects/procyon-core"

type HandlerAdapter interface {
	Supports(handler interface{}) bool
	Handle(handler interface{}, res HttpResponse, req HttpRequest) interface{}
}

type RequestMappingHandlerAdapter struct {
	typeConverterService core.TypeConverterService
	parameterResolvers   *HandlerMethodParameterResolvers
	returnValueHandlers  *HandlerMethodReturnValueHandlers
}

type RequestMappingHandlerAdapterOption func(adapter *RequestMappingHandlerAdapter)

func NewRequestMappingHandlerAdapter(service core.TypeConverterService) *RequestMappingHandlerAdapter {
	adapter := &RequestMappingHandlerAdapter{
		parameterResolvers:  getDefaultMethodParameterResolvers(service),
		returnValueHandlers: getDefaultReturnValueHandlers(),
	}
	return adapter
}

func getDefaultMethodParameterResolvers(service core.TypeConverterService) *HandlerMethodParameterResolvers {
	resolvers := NewHandlerMethodParameterResolvers()
	resolvers.AddMethodParameterResolver(NewContextMethodParameterResolver())
	resolvers.AddMethodParameterResolver(NewRequestMethodParameterResolver(service))
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

func (adapter *RequestMappingHandlerAdapter) Supports(handler interface{}) bool {
	if _, ok := handler.(HandlerMethod); ok {
		return true
	}
	return false
}

func (adapter *RequestMappingHandlerAdapter) Handle(handler interface{}, res HttpResponse, req HttpRequest) interface{} {
	return adapter.invokeHandler(handler.(HandlerMethod), res, req)
}

func (adapter *RequestMappingHandlerAdapter) invokeHandler(handler HandlerMethod, res HttpResponse, req HttpRequest) interface{} {
	arguments := adapter.getMethodArgumentValues(handler, req)
	results := handler.InvokeHandler(arguments)
	if len(results) > 0 {

	}
	return nil
}

func (adapter *RequestMappingHandlerAdapter) getMethodArgumentValues(handler HandlerMethod, req HttpRequest) []interface{} {
	if handler.GetHandlerParameterCount() == 0 {
		return nil
	}
	argumentValues := make([]interface{}, handler.GetHandlerParameterCount())
	for index, parameterType := range handler.GetHandlerParameterTypes() {
		if !adapter.parameterResolvers.SupportsParameter(parameterType) {
			panic("No suitable parameter resolver")
		}
		value, err := adapter.parameterResolvers.ResolveParameter(parameterType, req)
		if err != nil {
			panic(err)
		}
		argumentValues[index] = value
	}
	return argumentValues
}
