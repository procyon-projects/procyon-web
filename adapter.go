package web

import (
	"github.com/procyon-projects/procyon-core"
	"net/http"
)

type HandlerAdapter interface {
	Supports(handler interface{}, requestContext RequestContext) bool
	Handle(handler interface{}, requestContext RequestContext, res http.ResponseWriter, req *http.Request) interface{}
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

func (adapter *RequestMappingHandlerAdapter) Supports(handler interface{}, requestContext RequestContext) bool {
	if _, ok := handler.(HandlerMethod); ok {
		return true
	}
	return false
}

func (adapter *RequestMappingHandlerAdapter) Handle(handler interface{}, requestContext RequestContext, res http.ResponseWriter, req *http.Request) interface{} {
	return adapter.invokeHandler(handler.(HandlerMethod), requestContext, res, req)
}

func (adapter *RequestMappingHandlerAdapter) invokeHandler(handler HandlerMethod, requestContext RequestContext, res http.ResponseWriter, req *http.Request) interface{} {
	arguments := adapter.getMethodArgumentValues(handler, requestContext, req)
	results := handler.InvokeHandler(arguments)
	if len(results) > 0 {

	}
	return nil
}

func (adapter *RequestMappingHandlerAdapter) getMethodArgumentValues(handler HandlerMethod, requestContext RequestContext, req *http.Request) []interface{} {
	if handler.GetHandlerParameterCount() == 0 {
		return nil
	}
	argumentValues := make([]interface{}, handler.GetHandlerParameterCount())
	for index, parameterType := range handler.GetHandlerParameterTypes() {
		if !adapter.parameterResolvers.SupportsParameter(parameterType, requestContext, req) {
			panic("No suitable parameter resolver")
		}
		value, err := adapter.parameterResolvers.ResolveParameter(parameterType, requestContext, req)
		if err != nil {
			panic(err)
		}
		argumentValues[index] = value
	}
	return argumentValues
}
