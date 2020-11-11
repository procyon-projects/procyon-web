package web

import (
	"github.com/procyon-projects/procyon-core"
	"net/http"
)

type HandlerAdapter interface {
	Supports(handler *HandlerMethod, requestContext RequestContext) bool
	Handle(handler *HandlerMethod, requestContext RequestContext, res http.ResponseWriter, req *http.Request) interface{}
}

type RequestMappingHandlerAdapter struct {
	typeConverterService core.TypeConverterService
	/*parameterResolvers   *HandlerMethodParameterResolvers
	returnValueHandlers  *HandlerMethodReturnValueHandlers*/
}

type RequestMappingHandlerAdapterOption func(adapter *RequestMappingHandlerAdapter)

func NewRequestMappingHandlerAdapter(service core.TypeConverterService) *RequestMappingHandlerAdapter {
	adapter := &RequestMappingHandlerAdapter{
		//parameterResolvers:  getDefaultMethodParameterResolvers(service),
		//returnValueHandlers: getDefaultReturnValueHandlers(),
	}
	return adapter
}

/*
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
*/

func (adapter *RequestMappingHandlerAdapter) Supports(handler *HandlerMethod, requestContext RequestContext) bool {
	return true
}

func (adapter *RequestMappingHandlerAdapter) Handle(handler *HandlerMethod, requestContext RequestContext, res http.ResponseWriter, req *http.Request) interface{} {
	handler.method(requestContext)
	return nil
}
