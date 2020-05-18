package web

import "net/http"

type HttpMethod string

const (
	HttpMethodGet   HttpMethod = http.MethodGet
	HttpMethodHead  HttpMethod = http.MethodHead
	HttpMethodPost  HttpMethod = http.MethodPost
	HttpMethodPut   HttpMethod = http.MethodPut
	HttpMethodPatch HttpMethod = http.MethodPatch
)

type HandlerMethodRegistry struct {
	registerMap map[string][]*HandlerMethod
}

func NewHandlerMethodRegistry() *HandlerMethodRegistry {
	return &HandlerMethodRegistry{
		registerMap: make(map[string][]*HandlerMethod),
	}
}

func (registry *HandlerMethodRegistry) Register(handlerMethod *HandlerMethod) {
	registry.RegisterGroup("", handlerMethod)
}

func (registry *HandlerMethodRegistry) RegisterGroup(groupName string, handlerMethod *HandlerMethod) {
	if registry.registerMap[groupName] == nil {
		registry.registerMap[groupName] = make([]*HandlerMethod, 0)
	}
	registry.registerMap[groupName] = append(registry.registerMap[""], handlerMethod)
}
