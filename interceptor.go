package web

import "errors"

type HandlerInterceptor interface {
	HandleBefore(requestContext *WebRequestContext)
	HandleAfter(requestContext *WebRequestContext)
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
			requestContext.Next()
		}
	}
}

func (interceptor RecoveryInterceptor) HandleBefore(requestContext *WebRequestContext) {
	defer recoveryFunction(requestContext)
	requestContext.Next()
}

func (interceptor RecoveryInterceptor) HandleAfter(requestContext *WebRequestContext) {
	requestContext.Next()
}

func (interceptor RecoveryInterceptor) AfterCompletion(requestContext *WebRequestContext) {
	requestContext.Next()
}
