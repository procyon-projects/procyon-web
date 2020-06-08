package web

type HandlerMapping interface {
	GetHandlerChain(req HttpRequest) *HandlerChain
}
