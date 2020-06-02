package web

type HandlerAdapter interface {
	Supports(handler interface{}) bool
	Handle(handler interface{}, res HttpResponse, req HttpRequest) interface{}
}

type RequestMappingHandlerAdapter struct {
	parameterResolvers  *HandlerMethodParameterResolvers
	returnValueHandlers *HandlerMethodReturnValueHandlers
}

func NewRequestMappingHandlerAdapter() *RequestMappingHandlerAdapter {
	return &RequestMappingHandlerAdapter{
		parameterResolvers:  NewHandlerMethodParameterResolvers(),
		returnValueHandlers: NewHandlerMethodReturnValueHandlers(),
	}
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
