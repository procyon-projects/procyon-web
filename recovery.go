package web

import "errors"

func recoveryFunction(requestContext *WebRequestContext) {
	if r := recover(); r != nil {
		if !requestContext.completedFlow {
			switch val := r.(type) {
			case string:
				requestContext.err = errors.New(val)
			case error:
				requestContext.err = val
			default:
				requestContext.err = errors.New("unknown error")
			}
			requestContext.completedFlow = false
			requestContext.handlerIndex = requestContext.handlerChain.afterCompletionStartIndex - 1
		}
	}
}
