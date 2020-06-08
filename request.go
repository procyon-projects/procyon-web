package web

import (
	core "github.com/procyon-projects/procyon-core"
	"log"
	"net/http"
)

type RequestHandlerFunc = interface{}
type RequestHandlerOption func(handler *RequestHandlerInfo)

type HttpMethod string

const (
	HttpMethodGet    HttpMethod = http.MethodGet
	HttpMethodPost   HttpMethod = http.MethodPost
	HttpMethodPut    HttpMethod = http.MethodPut
	HttpMethodDelete HttpMethod = http.MethodDelete
	HttpMethodPatch  HttpMethod = http.MethodPatch
)

type RequestHandlerInfo struct {
	Paths       []string
	Methods     []HttpMethod
	HandlerFunc RequestHandlerFunc
}

func NewHandlerInfo(handler RequestHandlerFunc, options ...RequestHandlerOption) *RequestHandlerInfo {
	if handler == nil {
		log.Fatal("Handler must not be null")
	}
	typ := core.GetType(handler)
	if !core.IsFunc(typ) {
		log.Fatal("Handler must be function")
	}
	handlerMethod := &RequestHandlerInfo{
		HandlerFunc: handler,
	}
	for _, option := range options {
		option(handlerMethod)
	}
	if len(handlerMethod.Methods) == 0 {
		handlerMethod.Methods = []HttpMethod{HttpMethodGet}
	}
	return handlerMethod
}

func WithPath(paths ...string) RequestHandlerOption {
	return func(handler *RequestHandlerInfo) {
		handler.Paths = paths
	}
}

func WithMethod(methods ...HttpMethod) RequestHandlerOption {
	return func(handler *RequestHandlerInfo) {
		if len(methods) > 0 {
			handler.Methods = methods
		} else {
			handler.Methods = []HttpMethod{HttpMethodGet}
		}
	}
}

type HandlerInfoRegistry interface {
	Register(info ...*RequestHandlerInfo)
	RegisterGroup(prefix string, info ...*RequestHandlerInfo)
}

type SimpleHandlerInfoRegistry struct {
	infoRegistryMap map[string][]*RequestHandlerInfo
}

func NewSimpleHandlerInfoRegistry() *SimpleHandlerInfoRegistry {
	return &SimpleHandlerInfoRegistry{
		infoRegistryMap: make(map[string][]*RequestHandlerInfo),
	}
}

func (registry *SimpleHandlerInfoRegistry) Register(info ...*RequestHandlerInfo) {
	registry.RegisterGroup("<nil>", info...)
}

func (registry *SimpleHandlerInfoRegistry) RegisterGroup(prefix string, info ...*RequestHandlerInfo) {
	if len(info) == 0 {
		return
	}
	if prefix == "" {
		prefix = "<nil>"
	}
	if registry.infoRegistryMap[prefix] == nil {
		registry.infoRegistryMap[prefix] = make([]*RequestHandlerInfo, 0)
	}
	registry.infoRegistryMap[prefix] = append(registry.infoRegistryMap[prefix], info...)
}
