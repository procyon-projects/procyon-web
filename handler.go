package web

import (
	"github.com/codnect/goo"
)

type HandlerMethod struct {
	method RequestHandlerFunc
}

func NewSimpleHandlerMethod(method RequestHandlerFunc) *HandlerMethod {
	typ := goo.GetType(method)
	if !typ.IsFunction() {
		panic("Handler method is not function")
	}
	handlerMethod := &HandlerMethod{
		method,
	}
	return handlerMethod
}

func (handlerMethod HandlerMethod) InvokeHandler(args ...interface{}) []interface{} {
	return nil
}
