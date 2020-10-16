package web

import (
	"github.com/codnect/goo"
)

type ErrorHandlerFunc interface{}

type ErrorHandler struct {
	HandlerFunc ErrorHandlerFunc
	Errors      []error
}

func NewErrorHandler(handler ErrorHandlerFunc, errors ...error) ErrorHandler {
	if handler == nil {
		panic("Handler must not be null")
	}
	if errors == nil {
		panic("Error(s) must not be null")
	}
	typ := goo.GetType(handler)
	if !typ.IsFunction() {
		panic("Handler must be function")
	}
	handlerMethod := ErrorHandler{
		HandlerFunc: handler,
		Errors:      errors,
	}
	return handlerMethod
}

type ErrorHandlerRegistry interface {
	Register(handler ErrorHandler)
}

type ErrorHandlerAdviser interface {
	RegisterErrorHandlers(registry ErrorHandlerRegistry)
}

type RouterError struct {
	message string
}

func NewRouterError(message string) RouterError {
	return RouterError{message}
}

func (err RouterError) Error() string {
	return err.message
}

type NoHandlerFoundError struct {
	message string
}

func NewNoHandlerFoundError(message string) NoHandlerFoundError {
	return NoHandlerFoundError{message}
}

func (err NoHandlerFoundError) Error() string {
	return err.message
}

type NoHandlerParameterResolverError struct {
	message string
}

func NewNoHandlerParameterResolverError(message string) NoHandlerParameterResolverError {
	return NoHandlerParameterResolverError{message}
}

func (err NoHandlerParameterResolverError) Error() string {
	return err.message
}
