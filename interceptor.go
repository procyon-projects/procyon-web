package web

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
