package web

import (
	"github.com/codnect/goo"
)

type HandlerMethod interface {
	GetHandlerParameterCount() int
	GetHandlerReturnTypeCount() int
	GetHandlerParameterTypes() []HandlerMethodParameter
	GetHandlerReturnValues() []HandlerMethodReturnValue
}

type SimpleHandlerMethod struct {
	parameters   []HandlerMethodParameter
	returnValues []HandlerMethodReturnValue
}

func NewSimpleHandlerMethod(method interface{}) SimpleHandlerMethod {
	return SimpleHandlerMethod{
		make([]HandlerMethodParameter, 0),
		make([]HandlerMethodReturnValue, 0),
	}
}

func (handlerMethod SimpleHandlerMethod) GetHandlerParameterCount() int {
	return len(handlerMethod.parameters)
}

func (handlerMethod SimpleHandlerMethod) GetHandlerReturnTypeCount() int {
	return len(handlerMethod.returnValues)
}

func (handlerMethod SimpleHandlerMethod) GetHandlerParameterTypes() []HandlerMethodParameter {
	return handlerMethod.parameters
}

func (handlerMethod SimpleHandlerMethod) GetHandlerReturnValues() []HandlerMethodReturnValue {
	return handlerMethod.returnValues
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

type HandlerChain interface {
	GetHandler() interface{}
	GetHandlerInterceptors() []HandlerInterceptor
}

type HandlerExecutionChain struct {
	handler      interface{}
	interceptors []HandlerInterceptor
}

type HandlerExecutionChainOption func(chain *HandlerExecutionChain)

func NewHandlerExecutionChain(handler interface{}, options ...HandlerExecutionChainOption) HandlerExecutionChain {
	chain := &HandlerExecutionChain{
		handler: handler,
	}
	if len(options) == 0 {
		chain.interceptors = make([]HandlerInterceptor, 0)
	}
	for _, option := range options {
		option(chain)
	}
	return *chain
}

func WithInterceptors(interceptors []HandlerInterceptor) HandlerExecutionChainOption {
	return func(chain *HandlerExecutionChain) {
		chain.interceptors = interceptors
	}
}

func (chain HandlerExecutionChain) GetHandler() interface{} {
	return chain.handler
}

func (chain HandlerExecutionChain) GetHandlerInterceptors() []HandlerInterceptor {
	return chain.interceptors
}

func (chain HandlerExecutionChain) applyHandleBefore(res HttpResponse, req HttpRequest) {
	interceptors := chain.interceptors
	for _, interceptor := range interceptors {
		interceptor.HandleBefore(chain, res, req)
	}
}

func (chain HandlerExecutionChain) applyHandleAfter(res HttpResponse, req HttpRequest) {
	interceptors := chain.interceptors
	for _, interceptor := range interceptors {
		interceptor.HandleAfter(chain, res, req)
	}
}
