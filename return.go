package web

import core "github.com/Rollcomp/procyon-core"

type HandlerMethodReturnValueHandler interface {
	SupportsReturnType(returnValueType HandlerMethodReturnValue) bool
	HandleReturnValue(returnValues []interface{}, returnValueParameter HandlerMethodReturnValue, request HttpRequest) (interface{}, error)
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

func (h *HandlerMethodReturnValueHandlers) HandleReturnValue(returnValues []interface{}, returnValueParameter HandlerMethodReturnValue, request HttpRequest) (interface{}, error) {
	handler := h.findReturnValueHandler(returnValueParameter)
	if handler == nil {
		return nil, NewNoHandlerParameterResolver("Return value handler not found")
	}
	return handler.HandleReturnValue(returnValues, returnValueParameter, request)
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

type ResponseEntityReturnValueHandler struct {
}

func NewResponseBodyReturnValueHandler() ResponseEntityReturnValueHandler {
	return ResponseEntityReturnValueHandler{}
}

func (h ResponseEntityReturnValueHandler) SupportsReturnType(returnValueType HandlerMethodReturnValue) bool {
	if !returnValueType.HasType(core.GetType((*ResponseEntity)(nil))) {
		return false
	}
	if returnValueType.GetReturnTypeCount() == 2 {
		if returnValueType.HasErrorType() {
			return true
		}
		return false
	}
	return true
}

func (h ResponseEntityReturnValueHandler) HandleReturnValue(returnValues []interface{}, returnValueParameter HandlerMethodReturnValue, request HttpRequest) (interface{}, error) {
	/* TODO it will be completed */
	return nil, nil
}
