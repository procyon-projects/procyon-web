package web

import (
	"github.com/procyon-projects/goo"
	"net/http"
)

type RequestObjectCache struct {
	hasOnlyBody      bool
	bodyFieldIndex   int
	paramFieldIndex  int
	pathFieldIndex   int
	headerFieldIndex int
	fields           []goo.Field
}

type RequestHandlerObject interface{}
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

type RequestHandler struct {
	Path                  string
	Method                RequestMethod
	HandlerFunc           RequestHandlerFunction
	RequestObject         RequestHandlerObject
	requestObjectMetadata *RequestObjectMetadata
}

func newHandler(handler RequestHandlerFunction, method RequestMethod, options ...RequestHandlerOption) RequestHandler {
	if handler == nil {
		panic("Handler must not be null")
	}

	handlerType := goo.GetType(handler)
	if !handlerType.IsFunction() {
		panic("Handler must be function")
	}

	requestHandler := &RequestHandler{
		HandlerFunc: handler,
		Method:      method,
	}

	for _, option := range options {
		option(requestHandler)
	}

	if requestHandler.RequestObject != nil {
		requestHandler.requestObjectMetadata = ScanRequestObjectMetadata(requestHandler.RequestObject)
	}

	return *requestHandler
}

func Get(handler RequestHandlerFunction, options ...RequestHandlerOption) RequestHandler {
	return newHandler(handler, RequestMethodGet, options...)
}

func Post(handler RequestHandlerFunction, options ...RequestHandlerOption) RequestHandler {
	return newHandler(handler, RequestMethodPost, options...)
}

func Put(handler RequestHandlerFunction, options ...RequestHandlerOption) RequestHandler {
	return newHandler(handler, RequestMethodPut, options...)
}

func Delete(handler RequestHandlerFunction, options ...RequestHandlerOption) RequestHandler {
	return newHandler(handler, RequestMethodDelete, options...)
}

func Patch(handler RequestHandlerFunction, options ...RequestHandlerOption) RequestHandler {
	return newHandler(handler, RequestMethodPatch, options...)
}

func Options(handler RequestHandlerFunction, options ...RequestHandlerOption) RequestHandler {
	return newHandler(handler, RequestMethodOptions, options...)
}

func Head(handler RequestHandlerFunction, options ...RequestHandlerOption) RequestHandler {
	return newHandler(handler, RequestMethodHead, options...)
}

func RequestObject(requestObject RequestHandlerObject) RequestHandlerOption {
	return func(handler *RequestHandler) {
		handler.RequestObject = requestObject
	}
}

func Path(path string) RequestHandlerOption {
	return func(handler *RequestHandler) {
		handler.Path = path
	}
}

type HandlerRegistry interface {
	Register(info ...RequestHandler)
	RegisterGroup(prefix string, info ...RequestHandler)
}

type SimpleHandlerRegistry struct {
	registryMap map[string][]RequestHandler
}

func NewSimpleHandlerRegistry() SimpleHandlerRegistry {
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

	for _, handler := range info {
		if prefix+handler.Path == "" {
			panic("Specify a path or a prefix while registering a request handler")
		}
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
