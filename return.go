package web

import core "github.com/Rollcomp/procyon-core"

type HandlerMethodReturnValueHandler interface {
	SupportsReturnType(returnValueType HandlerMethodReturnValue) bool
	HandleReturnValue(returnValue interface{}, returnValueParameter HandlerMethodReturnValue, request HttpRequest) (interface{}, error)
}

type HandlerMethodReturnValueHandlers struct {
	returnValueHandlers []HandlerMethodReturnValueHandler
}

func NewHandlerMethodReturnValueHandlers() *HandlerMethodReturnValueHandlers {
	return &HandlerMethodReturnValueHandlers{
		make([]HandlerMethodReturnValueHandler, 0),
	}
}

func (h *HandlerMethodReturnValueHandlers) SupportsReturnType(returnValueType HandlerMethodReturnValue) bool {
	for _, handler := range h.returnValueHandlers {
		if handler.SupportsReturnType(returnValueType) {
			return true
		}
	}
	return false
}

func (h *HandlerMethodReturnValueHandlers) HandleReturnValue(returnValue interface{}, returnValueParameter HandlerMethodReturnValue, request HttpRequest) (interface{}, error) {
	handler := h.findReturnValueHandler(returnValueParameter)
	if handler == nil {
		return nil, NewNoHandlerParameterResolver("Return value handler not found")
	}
	return handler.HandleReturnValue(returnValue, returnValueParameter, request)
}

func (h *HandlerMethodReturnValueHandlers) findReturnValueHandler(returnValueParameter HandlerMethodReturnValue) HandlerMethodReturnValueHandler {
	for _, handler := range h.returnValueHandlers {
		if handler.SupportsReturnType(returnValueParameter) {
			return handler
		}
	}
	return nil
}

func (h *HandlerMethodReturnValueHandlers) AddMethodReturnValueHandler(handlers ...HandlerMethodReturnValueHandler) {
	h.returnValueHandlers = append(h.returnValueHandlers, handlers...)
}

type ResponseBodyReturnValueHandler struct {
}

func NewResponseBodyReturnValueHandler() ResponseBodyReturnValueHandler {
	return ResponseBodyReturnValueHandler{}
}

func (h ResponseBodyReturnValueHandler) SupportsReturnType(returnValueType HandlerMethodReturnValue) bool {
	if returnValueType.GetReturnTypeCount() == 2 && returnValueType.HasErrorType() &&
		returnValueType.HasType(core.GetType((*ResponseBody)(nil))) {
		return true
	}
	return false
}

func (h ResponseBodyReturnValueHandler) HandleReturnValue(returnValue interface{}, returnValueParameter HandlerMethodReturnValue, request HttpRequest) (interface{}, error) {
	/* TODO it will be completed */
	return nil, nil
}
