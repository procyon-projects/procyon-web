package web

import "errors"

func recoveryFunction(requestContext *WebRequestContext) {
	r := recover()
	if r != nil {
		requestContext.crashed = true

		switch val := r.(type) {
		case string:
			requestContext.err = errors.New(val)
		case error:
			requestContext.err = val
		default:
			requestContext.err = errors.New("unknown error")
		}

		if requestContext.handlerChain != nil {
			requestContext.handlerIndex = requestContext.handlerChain.afterCompletionStartIndex
			requestContext.invokeHandlers(true)
		}
	}
}
