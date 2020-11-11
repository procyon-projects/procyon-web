package web

/*
type HandlerMethodReturnValueHandler interface {
	SupportsReturnType(returnValueType HandlerMethodReturnValue) bool
	HandleReturnValue(returnValues []interface{}, returnValueParameter HandlerMethodReturnValue, request *http.Request) (interface{}, error)
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

func (h *HandlerMethodReturnValueHandlers) HandleReturnValue(returnValues []interface{}, returnValueParameter HandlerMethodReturnValue, request *http.Request) (interface{}, error) {
	handler := h.findReturnValueHandler(returnValueParameter)
	if handler == nil {
		return nil, NewNoHandlerParameterResolverError("Return value handler not found")
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

func NewResponseEntityReturnValueHandler() ResponseEntityReturnValueHandler {
	return ResponseEntityReturnValueHandler{}
}

func (h ResponseEntityReturnValueHandler) SupportsReturnType(returnValueType HandlerMethodReturnValue) bool {
	if returnValueType.GetReturnType().Equals(goo.GetType((*ResponseEntity)(nil))) {
		return true
	}
	return false
}

func (h ResponseEntityReturnValueHandler) HandleReturnValue(returnValues []interface{}, returnValueParameter HandlerMethodReturnValue, request *http.Request) (interface{}, error) {
	return nil, nil
}

type ErrorReturnValueHandler struct {
}

func NewErrorReturnValueHandler() ErrorReturnValueHandler {
	return ErrorReturnValueHandler{}
}

func (h ErrorReturnValueHandler) SupportsReturnType(returnValueType HandlerMethodReturnValue) bool {
	if returnValueType.GetReturnType().Equals(goo.GetType((error)(nil))) {
		return true
	}
	return false
}

func (h ErrorReturnValueHandler) HandleReturnValue(returnValues []interface{}, returnValueParameter HandlerMethodReturnValue, request *http.Request) (interface{}, error) {
	return nil, nil
}
*/
