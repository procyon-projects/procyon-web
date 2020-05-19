package web

import (
	"log"
	"net/http"
)
import "github.com/codnect/go-reflect"

type HandlerFunc = interface{}
type HttpMethod string

const (
	HttpMethodGet    HttpMethod = http.MethodGet
	HttpMethodPost   HttpMethod = http.MethodPost
	HttpMethodPut    HttpMethod = http.MethodPut
	HttpMethodDelete HttpMethod = http.MethodDelete
	HttpMethodPatch  HttpMethod = http.MethodPatch
)

type HandlerMethod struct {
	Path        string
	Method      HttpMethod
	HandlerFunc HandlerFunc
}

func newHandlerMethod(path string, method HttpMethod, handler HandlerFunc) *HandlerMethod {
	if handler == nil {
		log.Fatal("Handler must not be null")
	}
	typ := reflect.GetType(handler)
	if !typ.IsFunction() {
		log.Fatal("Handler must be function")
	}
	return &HandlerMethod{
		Path:        path,
		Method:      method,
		HandlerFunc: handler,
	}
}

func WithGet(path string, handler HandlerFunc) *HandlerMethod {
	return newHandlerMethod(path, HttpMethodGet, handler)
}

func WithPost(path string, handler HandlerFunc) *HandlerMethod {
	return newHandlerMethod(path, HttpMethodPost, handler)
}

func WithDelete(path string, handler HandlerFunc) *HandlerMethod {
	return newHandlerMethod(path, HttpMethodDelete, handler)
}

func WithPut(path string, handler HandlerFunc) *HandlerMethod {
	return newHandlerMethod(path, HttpMethodPut, handler)
}

func WithPatch(path string, handler HandlerFunc) *HandlerMethod {
	return newHandlerMethod(path, HttpMethodPatch, handler)
}

type HandlerMethodRegistry struct {
	registerMap map[string][]*HandlerMethod
}

func NewHandlerMethodRegistry() *HandlerMethodRegistry {
	return &HandlerMethodRegistry{
		registerMap: make(map[string][]*HandlerMethod),
	}
}

func (registry *HandlerMethodRegistry) Register(handlerMethod ...*HandlerMethod) {
	registry.RegisterGroup("", handlerMethod...)
}

func (registry *HandlerMethodRegistry) RegisterGroup(groupName string, handlerMethod ...*HandlerMethod) {
	if handlerMethod == nil {
		return
	}
	if registry.registerMap[groupName] == nil {
		registry.registerMap[groupName] = make([]*HandlerMethod, 0)
	}
	registry.registerMap[groupName] = append(registry.registerMap[groupName], handlerMethod...)
}
