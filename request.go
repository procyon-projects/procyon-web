package web

import (
	"github.com/codnect/goo"
	"net/http"
)

type RequestHandlerFunction = func(context *WebRequestContext)
type RequestHandlerOption func(handler *RequestHandler)

type RequestMethod string

const (
	unknownMethod        RequestMethod = "[unknown-method]"
	RequestMethodGet     RequestMethod = http.MethodGet
	RequestMethodPost    RequestMethod = http.MethodPost
	RequestMethodPut     RequestMethod = http.MethodPut
	RequestMethodDelete  RequestMethod = http.MethodDelete
	RequestMethodPatch   RequestMethod = http.MethodPatch
	RequestMethodOptions RequestMethod = http.MethodOptions
	RequestMethodHead    RequestMethod = http.MethodHead
)

func GetRequestMethod(method string) RequestMethod {
	switch method {
	case http.MethodGet:
		return RequestMethodGet
	case http.MethodPost:
		return RequestMethodPost
	case http.MethodPut:
		return RequestMethodPut
	case http.MethodDelete:
		return RequestMethodDelete
	case http.MethodPatch:
		return RequestMethodPatch
	case http.MethodOptions:
		return RequestMethodOptions
	case http.MethodHead:
		return RequestMethodHead
	}
	return unknownMethod
}

type RequestHandler struct {
	Path        string
	Methods     []RequestMethod
	HandlerFunc RequestHandlerFunction
}

func NewHandler(handler RequestHandlerFunction, options ...RequestHandlerOption) RequestHandler {
	if handler == nil {
		panic("Handler must not be null")
	}
	typ := goo.GetType(handler)
	if !typ.IsFunction() {
		panic("Handler must be function")
	}
	handlerMethod := &RequestHandler{
		HandlerFunc: handler,
	}
	for _, option := range options {
		option(handlerMethod)
	}
	if len(handlerMethod.Methods) == 0 {
		handlerMethod.Methods = []RequestMethod{RequestMethodGet}
	}
	return *handlerMethod
}

func WithPath(path string) RequestHandlerOption {
	return func(handler *RequestHandler) {
		handler.Path = path
	}
}

func WithMethod(methods ...RequestMethod) RequestHandlerOption {
	return func(handler *RequestHandler) {
		if len(methods) > 0 {
			handler.Methods = methods
		} else {
			handler.Methods = []RequestMethod{RequestMethodGet}
		}
	}
}

type HandlerRegistry interface {
	Register(info ...RequestHandler)
	RegisterGroup(prefix string, info ...RequestHandler)
}

type SimpleHandlerRegistry struct {
	registryMap map[string][]RequestHandler
}

func newSimpleHandlerRegistry() SimpleHandlerRegistry {
	return SimpleHandlerRegistry{
		registryMap: make(map[string][]RequestHandler),
	}
}

func (registry SimpleHandlerRegistry) Register(info ...RequestHandler) {
	registry.RegisterGroup("", info...)
}

func (registry SimpleHandlerRegistry) RegisterGroup(prefix string, info ...RequestHandler) {
	if len(info) == 0 {
		return
	}
	if registry.registryMap[prefix] == nil {
		registry.registryMap[prefix] = make([]RequestHandler, 0)
	}
	registry.registryMap[prefix] = append(registry.registryMap[prefix], info...)
}

func (registry SimpleHandlerRegistry) clear() {
	for key := range registry.registryMap {
		delete(registry.registryMap, key)
	}
}

func (registry SimpleHandlerRegistry) getRegistryMap() map[string][]RequestHandler {
	return registry.registryMap
}
