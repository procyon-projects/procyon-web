package web

import (
	"github.com/codnect/goo"
)

type HandlerMethod struct {
	peaName      string
	parameters   []HandlerMethodParameter
	returnValues []HandlerMethodReturnValue
}

func NewHandlerMethod(peaName string, method interface{}) HandlerMethod {
	return HandlerMethod{
		peaName,
		make([]HandlerMethodParameter, 0),
		make([]HandlerMethodReturnValue, 0),
	}
}

func (m HandlerMethod) GetParameterCount() int {
	return len(m.parameters)
}

func (m HandlerMethod) GetReturnTypeCount() int {
	return len(m.returnValues)
}

func (m HandlerMethod) GetParameterTypes() []HandlerMethodParameter {
	return m.parameters
}

func (m HandlerMethod) GetReturnValues() []HandlerMethodReturnValue {
	return m.returnValues
}

type HandlerMethodParameter struct {
	typ goo.Type
}

func NewHandlerMethodParameter(typ goo.Type) HandlerMethodParameter {
	return HandlerMethodParameter{
		typ,
	}
}

func (r HandlerMethodParameter) GetType() goo.Type {
	return r.typ
}

type HandlerMethodReturnValue struct {
	typ goo.Type
}

func NewHandlerMethodReturnValue(typ goo.Type) HandlerMethodReturnValue {
	return HandlerMethodReturnValue{
		typ,
	}
}

func (r HandlerMethodReturnValue) GetType() goo.Type {
	return r.typ
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
	interceptors := chain.interceptors
	for _, interceptor := range interceptors {
		interceptor.HandleBefore(chain, res, req)
	}
}

func (chain *HandlerChain) applyHandleAfter(res HttpResponse, req HttpRequest) {
	interceptors := chain.interceptors
	for _, interceptor := range interceptors {
		interceptor.HandleAfter(chain, res, req)
	}
}
