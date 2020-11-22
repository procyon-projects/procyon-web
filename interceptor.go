package web

import (
	"errors"
)

type HandlerInterceptor func(requestContext *WebRequestContext)

type HandlerInterceptorBefore interface {
	HandleBefore(requestContext *WebRequestContext)
}

type HandlerInterceptorAfter interface {
	HandleAfter(requestContext *WebRequestContext)
}

type HandlerInterceptorAfterCompletion interface {
	AfterCompletion(requestContext *WebRequestContext)
}

type RecoveryInterceptor struct {
}

func NewRecoveryInterceptor() RecoveryInterceptor {
	return RecoveryInterceptor{}
}

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
			requestContext.inMainHandler = false
			requestContext.handlerIndex = requestContext.handlerChain.afterCompletionStartIndex - 1
			requestContext.internalNext()
		}
	}
}

func (interceptor RecoveryInterceptor) HandleBefore(requestContext *WebRequestContext) {
	defer recoveryFunction(requestContext)
	requestContext.Next()
}
