package web

type HandlerFunction func(requestContext *WebRequestContext)

type HandlerChain struct {
	handler                   RequestHandlerFunction
	interceptors              []HandlerInterceptor
	allHandlers               []HandlerFunction
	handlerIndex              int
	afterStartIndex           int
	afterCompletionStartIndex int
	handlerEndIndex           int
	pathVariables             []string
}

func NewHandlerChain(fun RequestHandlerFunction, interceptors []HandlerInterceptor) *HandlerChain {
	chain := &HandlerChain{
		fun,
		interceptors,
		make([]HandlerFunction, 0),
		0,
		0,
		0,
		0,
		nil,
	}
	if len(chain.interceptors) == 0 {
		chain.interceptors = make([]HandlerInterceptor, 0)
	}
	for _, interceptor := range chain.interceptors {
		chain.allHandlers = append(chain.allHandlers, interceptor.HandleBefore)
	}
	chain.handlerIndex = len(chain.allHandlers)
	chain.allHandlers = append(chain.allHandlers, chain.handler)
	chain.afterStartIndex = len(chain.allHandlers)
	for index := len(interceptors) - 1; index >= 0; index-- {
		chain.allHandlers = append(chain.allHandlers, interceptors[index].HandleAfter)
	}
	chain.afterCompletionStartIndex = len(chain.allHandlers)
	for index := len(interceptors) - 1; index >= 0; index-- {
		chain.allHandlers = append(chain.allHandlers, interceptors[index].AfterCompletion)
	}
	chain.handlerEndIndex = len(chain.allHandlers) - 1
	return chain
}
