package web

type HandlerInterceptor interface {
	HandleBefore(handler interface{}, res HttpResponse, req HttpRequest)
	HandleAfter(handler interface{}, res HttpResponse, req HttpRequest)
}
