package web

import (
	"github.com/codnect/goo"
)

type HandlerMethod interface {
	GetHandlerParameterCount() int
	GetHandlerReturnTypeCount() int
	GetHandlerParameterTypes() []HandlerMethodParameter
	GetHandlerReturnValues() []HandlerMethodReturnValue
	InvokeHandler(args []interface{}) []interface{}
}

type SimpleHandlerMethod struct {
	method       interface{}
	methodType   goo.Function
	parameters   []HandlerMethodParameter
	returnValues []HandlerMethodReturnValue
}

func NewSimpleHandlerMethod(method interface{}) SimpleHandlerMethod {
	typ := goo.GetType(method)
	if !typ.IsFunction() {
		panic("Handler method is not function")
	}
	handlerMethod := &SimpleHandlerMethod{
		method,
		typ.ToFunctionType(),
		make([]HandlerMethodParameter, 0),
		make([]HandlerMethodReturnValue, 0),
	}
	handlerMethod.initHandlerMethodParameters()
	handlerMethod.initHandlerMethodReturnValues()
	return *handlerMethod
}

func (handlerMethod *SimpleHandlerMethod) initHandlerMethodParameters() {
	parameterTypes := handlerMethod.methodType.GetFunctionParameterTypes()
	for index, parameterType := range parameterTypes {
		handlerMethod.parameters = append(handlerMethod.parameters, NewHandlerMethodParameter(index, parameterType, handlerMethod.methodType))
	}
}

func (handlerMethod *SimpleHandlerMethod) initHandlerMethodReturnValues() {
	returnTypes := handlerMethod.methodType.GetFunctionReturnTypes()
	for index, returnType := range returnTypes {
		handlerMethod.returnValues = append(handlerMethod.returnValues, NewHandlerMethodReturnValue(index, returnType, handlerMethod.methodType))
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

func (handlerMethod SimpleHandlerMethod) InvokeHandler(args []interface{}) []interface{} {
	return handlerMethod.methodType.Call(args)
}

type HandlerMethodParameter struct {
	hashCode      int
	index         int
	parameterType goo.Type
	function      goo.Function
}

func NewHandlerMethodParameter(index int, parameterType goo.Type, function goo.Function) HandlerMethodParameter {
	hashCode := hashCodeForString(parameterType.GetFullName())
	return HandlerMethodParameter{
		hashCode,
		index,
		parameterType,
		function,
	}
}

func (handlerMethodParameter HandlerMethodParameter) GetParameterIndex() int {
	return handlerMethodParameter.index
}

func (handlerMethodParameter HandlerMethodParameter) GetParameterType() goo.Type {
	return handlerMethodParameter.parameterType
}

func (handlerMethodParameter HandlerMethodParameter) GetFunction() goo.Function {
	return handlerMethodParameter.function
}

func (handlerMethodParameter HandlerMethodParameter) HashCode() int {
	return handlerMethodParameter.hashCode
}

type HandlerMethodReturnValue struct {
	hashCode   int
	index      int
	returnType goo.Type
	function   goo.Function
}

func NewHandlerMethodReturnValue(index int, returnType goo.Type, function goo.Function) HandlerMethodReturnValue {
	hashCode := hashCodeForString(returnType.GetFullName())
	return HandlerMethodReturnValue{
		hashCode,
		index,
		returnType,
		function,
	}
}

func (handlerMethodReturnValue HandlerMethodReturnValue) GetReturnValueIndex() int {
	return handlerMethodReturnValue.index
}

func (handlerMethodReturnValue HandlerMethodReturnValue) GetReturnType() goo.Type {
	return handlerMethodReturnValue.returnType
}

func (handlerMethodReturnValue HandlerMethodReturnValue) GetFunction() goo.Function {
	return handlerMethodReturnValue.function
}

func (handlerMethodReturnValue HandlerMethodReturnValue) HashCode() int {
	return handlerMethodReturnValue.hashCode
}
