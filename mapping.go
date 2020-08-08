package web

import (
	"errors"
	"sync"
)

type HandlerMapping interface {
	GetHandlerChain(req HttpRequest) *HandlerChain
}

type mappingRegistry interface {
	Register(mapping string, handler interface{}, fun interface{}) error
	GetMappings() map[string]HandlerMethod
	FindMappingByUrl(path string) (string, error)
}

type defaultMappingRegistry struct {
	mappingLookup map[string]HandlerMethod
	urlLookup     map[string]string
	mu            sync.RWMutex
}

func newDefaultMappingRegistry() defaultMappingRegistry {
	return defaultMappingRegistry{
		mappingLookup: make(map[string]HandlerMethod),
		urlLookup:     make(map[string]string),
		mu:            sync.RWMutex{},
	}
}

func (registry defaultMappingRegistry) Register(mapping string, handler interface{}, fun interface{}) error {
	registry.mu.Lock()
	if _, ok := registry.mappingLookup[mapping]; ok {
		registry.mu.Unlock()
		return errors.New("ambiguous handler mapping. there is already an mapping :" + mapping)
	}
	registry.mappingLookup[mapping] = NewHandlerMethod(handler)
	registry.mu.Unlock()
	return nil
}

func (registry defaultMappingRegistry) GetMappings() map[string]HandlerMethod {
	return registry.mappingLookup
}

func (registry defaultMappingRegistry) FindMappingByUrl(path string) (string, error) {
	if result, ok := registry.urlLookup[path]; ok {
		return result, nil
	}
	return "", errors.New("not found matching")
}

type RequestHandlerMapping struct {
	mappingRegistry mappingRegistry
}

func NewRequestHandlerMapping() RequestHandlerMapping {
	return RequestHandlerMapping{
		mappingRegistry: newDefaultMappingRegistry(),
	}
}

func (requestMapping RequestHandlerMapping) RegisterHandlerMethod(mapping string, handler interface{}) error {
	return nil
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
