package web

import (
	"errors"
	"sync"
)

type HandlerMapping interface {
	GetHandlerChain(req HttpRequest) *HandlerChain
}

type MappingRegistry interface {
	Register(handlerName string, mapping RequestMappingInfo, fun RequestHandlerFunc) error
	GetMappings() map[int]HandlerMethod
	FindMappingByUrl(path string) (string, error)
}

type defaultMappingRegistry struct {
	nameLookup           map[string]HandlerMethod
	mappingLookup        map[int]RequestMappingInfo
	mappingHandlerLookup map[int]HandlerMethod
	mappingUrlLookup     map[string]RequestMappingInfo
	mu                   sync.RWMutex
}

func newDefaultMappingRegistry() defaultMappingRegistry {
	return defaultMappingRegistry{
		nameLookup:           make(map[string]HandlerMethod),
		mappingLookup:        make(map[int]RequestMappingInfo),
		mappingHandlerLookup: make(map[int]HandlerMethod),
		mappingUrlLookup:     make(map[string]RequestMappingInfo),
		mu:                   sync.RWMutex{},
	}
}

func (registry defaultMappingRegistry) Register(handlerName string, mapping RequestMappingInfo, fun RequestHandlerFunc) error {
	registry.mu.Lock()
	mappingHashCode := mapping.hashCode()
	if _, ok := registry.mappingLookup[mappingHashCode]; ok {
		registry.mu.Unlock()
		return errors.New("ambiguous handler mapping. there is already an mapping :" + handlerName)
	}
	registry.mappingLookup[mappingHashCode] = mapping
	registry.mappingHandlerLookup[mappingHashCode] = NewHandlerMethod(handlerName, fun)
	registry.findPurePaths(mapping)
	registry.mu.Unlock()
	return nil
}

func (registry defaultMappingRegistry) GetMappings() map[int]HandlerMethod {
	return nil
}

func (registry defaultMappingRegistry) FindMappingByUrl(path string) (string, error) {
	if _, ok := registry.mappingUrlLookup[path]; ok {
		return "", nil
	}
	return "", errors.New("not found matching")
}

func (registry defaultMappingRegistry) findPurePaths(mapping RequestMappingInfo) {
	patterns := mapping.getPatternRequestMatcher().patterns
	for _, pattern := range patterns {
		if !registry.isPatternPath(pattern) {
			registry.mappingUrlLookup[pattern] = mapping
		}
	}
}

func (registry defaultMappingRegistry) isPatternPath(path string) bool {
	pathVariable := false
	for _, character := range path {
		if character == '*' || character == '?' {
			return true
		}
		if character == '{' {
			pathVariable = true
		} else if pathVariable && character == '}' {
			return true
		}
	}
	return false
}

type RequestHandlerMapping struct {
	mappingRegistry MappingRegistry
	mu              sync.Mutex
}

func NewRequestHandlerMapping() RequestHandlerMapping {
	return RequestHandlerMapping{
		mappingRegistry: newDefaultMappingRegistry(),
	}
}

func (requestMapping RequestHandlerMapping) RegisterHandlerMethod(handlerName string, mapping RequestMappingInfo, fun RequestHandlerFunc) {
	requestMapping.mappingRegistry.Register(handlerName, mapping, fun)
}

func (requestMapping RequestHandlerMapping) GetHandlerChain(req HttpRequest) *HandlerChain {
	return nil
}

func (requestMapping RequestHandlerMapping) getHandler(req HttpRequest) {
	requestMapping.findHandlerMethod("", req)
}

func (requestMapping RequestHandlerMapping) findHandlerMethod(path string, req HttpRequest) {
	matches := make([]match, 0)
	directPathMapping, err := requestMapping.mappingRegistry.FindMappingByUrl(path)
	if err == nil {
		requestMapping.addMatches(matches, directPathMapping, req)
	}
	if len(matches) > 0 {

	}
}

func (requestMapping RequestHandlerMapping) addMatches(matches []match, mapping string, req HttpRequest) {

}

type match struct {
	mapping       string
	handlerMethod HandlerMethod
}

func newMatch(mapping string, method HandlerMethod) match {
	return match{
		mapping,
		method,
	}
}
