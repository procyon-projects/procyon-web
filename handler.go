package web

import core "github.com/Rollcomp/procyon-core"

type HandlerMethod struct {
	parameters []HandlerMethodParameter
}

type HandlerMethodParameter struct {
	typ *core.Type
}

type HandlerMethodReturnValue struct {
	typ []*core.Type
}

type HandlerChain struct {
	handler      interface{}
	interceptors []HandlerInterceptor
}

type HandlerChainOption func(chain *HandlerChain)

func NewHandlerExecutionChain(handler interface{}, options ...HandlerChainOption) *HandlerChain {
	chain := &HandlerChain{
		handler: handler,
	}
	if len(options) == 0 {
		chain.interceptors = make([]HandlerInterceptor, 0)
	}
	for _, option := range options {
		option(chain)
	}
	return chain
}

func WithInterceptors(interceptors []HandlerInterceptor) HandlerChainOption {
	return func(chain *HandlerChain) {
		chain.interceptors = interceptors
	}
}

func (chain *HandlerChain) getHandler() interface{} {
	return chain.handler
}

func (chain *HandlerChain) getInterceptors() []HandlerInterceptor {
	return chain.interceptors
}

func (chain *HandlerChain) applyHandleBefore(res HttpResponse, req HttpRequest) {
	for _, interceptor := range chain.interceptors {
		interceptor.HandleBefore(chain, res, req)
	}
}

func (chain *HandlerChain) applyHandleAfter(res HttpResponse, req HttpRequest) {
	for _, interceptor := range chain.interceptors {
		interceptor.HandleAfter(chain, res, req)
	}
}

type HandlerInterceptor interface {
	HandleBefore(handler interface{}, res HttpResponse, req HttpRequest)
	HandleAfter(handler interface{}, res HttpResponse, req HttpRequest)
}

type HandlerMapping interface {
	GetHandlerChain(req HttpRequest) *HandlerChain
}
