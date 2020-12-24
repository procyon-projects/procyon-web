package web

type HandlerFunction func(requestContext *WebRequestContext)

type HandlerChain struct {
	handler                   RequestHandlerFunction
	handlers                  []HandlerFunction
	handlerIndex              int
	afterStartIndex           int
	afterCompletionStartIndex int
	handlerEndIndex           int
	pathVariables             []string
	requestObjectMetadata     *RequestObjectMetadata
}

func NewHandlerChain(fun RequestHandlerFunction, interceptorRegistry HandlerInterceptorRegistry, metadata *RequestObjectMetadata) *HandlerChain {
	chain := &HandlerChain{
		fun,
		make([]HandlerFunction, 0),
		0,
		0,
		0,
		0,
		nil,
		metadata,
	}

	if interceptorRegistry != nil {
		for _, interceptor := range interceptorRegistry.GetHandlerBeforeInterceptors() {
			chain.handlers = append(chain.handlers, HandlerFunction(interceptor))
		}

		chain.handlerIndex = len(chain.handlers)
		chain.handlers = append(chain.handlers, chain.handler)

		chain.afterStartIndex = len(chain.handlers)
		for _, interceptor := range interceptorRegistry.GetHandlerAfterInterceptors() {
			chain.handlers = append(chain.handlers, HandlerFunction(interceptor))
		}

		chain.afterCompletionStartIndex = len(chain.handlers)
		for _, interceptor := range interceptorRegistry.GetHandlerAfterCompletionInterceptors() {
			chain.handlers = append(chain.handlers, HandlerFunction(interceptor))
		}

		chain.handlerEndIndex = len(chain.handlers) - 1

	} else {
		chain.handlers = append(chain.handlers)
		chain.handlerIndex = len(chain.handlers)

		chain.handlers = append(chain.handlers, chain.handler)
		chain.afterStartIndex = len(chain.handlers)

		chain.afterCompletionStartIndex = len(chain.handlers)
		chain.handlerEndIndex = len(chain.handlers) - 1
	}
	return chain
}

func (chain *HandlerChain) updatePathVariableMetadata(pathVariableIndex int, pathVariableName string) {
	if chain.requestObjectMetadata == nil {
		return
	}

	pathVariableMetadata := chain.requestObjectMetadata.pathMetadata

	for variableName, metadata := range pathVariableMetadata.pathVariableMap {
		if pathVariableName == variableName {
			metadata.extra = pathVariableIndex
			break
		}
	}

}
