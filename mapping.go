package web

import (
	"errors"
	"sync"
)

type RequestMapping struct {
	name                  string
	methodRequestMatcher  MethodRequestMatcher
	paramsRequestMatcher  ParametersRequestMatcher
	patternRequestMatcher PatternRequestMatcher
}

func NewRequestMapping(name string,
	methodRequestMatcher MethodRequestMatcher,
	paramsRequestMatcher ParametersRequestMatcher,
	patternRequestMatcher PatternRequestMatcher) RequestMapping {
	return RequestMapping{
		name,
		methodRequestMatcher,
		paramsRequestMatcher,
		patternRequestMatcher,
	}
}

func (mappingInfo RequestMapping) getMethodRequestMatcher() MethodRequestMatcher {
	return mappingInfo.methodRequestMatcher
}

func (mappingInfo RequestMapping) getParametersRequestMatcher() ParametersRequestMatcher {
	return mappingInfo.paramsRequestMatcher
}

func (mappingInfo RequestMapping) getPatternRequestMatcher() PatternRequestMatcher {
	return mappingInfo.patternRequestMatcher
}

func (mappingInfo RequestMapping) MatchRequest(req HttpRequest) interface{} {
	method := mappingInfo.methodRequestMatcher.MatchRequest(req)
	if method == nil {
		return nil
	} else {
		params := mappingInfo.paramsRequestMatcher.MatchRequest(req)
		if params == nil {
			return nil
		} else {
			pattern := mappingInfo.patternRequestMatcher.MatchRequest(req)
			if pattern == nil {
				return nil
			}
		}
	}
	return mappingInfo
}

func (mappingInfo RequestMapping) hashCode() int {
	return 31*mappingInfo.patternRequestMatcher.hashCode() +
		mappingInfo.methodRequestMatcher.hashCode() +
		mappingInfo.paramsRequestMatcher.hashCode()
}

type MappingRegistry interface {
	Register(handlerName string, mapping interface{}, fun RequestHandlerFunc) error
	GetMappings() map[interface{}]HandlerMethod
	FindMappingsByUrl(path string) ([]interface{}, error)
}

type RequestMappingRegistry struct {
	mappingHashcodeLookup map[int]bool
	mappingLookup         map[interface{}]HandlerMethod
	mappingUrlLookup      map[string][]interface{}
	mu                    sync.RWMutex
}

func NewRequestMappingRegistry() RequestMappingRegistry {
	return RequestMappingRegistry{
		mappingHashcodeLookup: make(map[int]bool),
		mappingLookup:         make(map[interface{}]HandlerMethod),
		mappingUrlLookup:      make(map[string][]interface{}),
		mu:                    sync.RWMutex{},
	}
}

func (registry RequestMappingRegistry) Register(handlerName string, mapping interface{}, fun RequestHandlerFunc) error {
	registry.mu.Lock()
	requestMapping := mapping.(RequestMapping)
	if _, ok := registry.mappingHashcodeLookup[requestMapping.hashCode()]; ok {
		registry.mu.Unlock()
		return errors.New("ambiguous handler mapping. there is already an mapping :" + handlerName)
	}
	registry.mappingHashcodeLookup[requestMapping.hashCode()] = true
	registry.mappingLookup[mapping] = NewSimpleHandlerMethod(fun)
	registry.findPurePaths(requestMapping)
	registry.mu.Unlock()
	return nil
}

func (registry RequestMappingRegistry) GetMappings() map[interface{}]HandlerMethod {
	return registry.mappingLookup
}

func (registry RequestMappingRegistry) FindMappingsByUrl(path string) ([]interface{}, error) {
	if mappings, ok := registry.mappingUrlLookup[path]; ok {
		return mappings, nil
	}
	return nil, errors.New("not found matching")
}

func (registry RequestMappingRegistry) findPurePaths(mapping RequestMapping) {
	patterns := mapping.getPatternRequestMatcher().patterns
	for _, pattern := range patterns {
		if !registry.isPatternPath(pattern) {
			registry.mappingUrlLookup[pattern] = append(registry.mappingUrlLookup[pattern], mapping)
		}
	}
}

func (registry RequestMappingRegistry) isPatternPath(path string) bool {
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

type HandlerMapping interface {
	RegisterHandlerMethod(handlerName string, mapping interface{}, fun RequestHandlerFunc)
	GetHandlerChain(req HttpRequest) HandlerChain
}

type RequestHandlerMapping struct {
	mappingRegistry MappingRegistry
	mu              sync.Mutex
}

func NewRequestHandlerMapping() RequestHandlerMapping {
	return RequestHandlerMapping{
		mappingRegistry: NewRequestMappingRegistry(),
	}
}

func (requestMapping RequestHandlerMapping) RegisterHandlerMethod(handlerName string, mapping interface{}, fun RequestHandlerFunc) {
	requestMapping.mappingRegistry.Register(handlerName, mapping, fun)
}

func (requestMapping RequestHandlerMapping) GetHandlerChain(req HttpRequest) HandlerChain {
	handler := requestMapping.lookupHandlerMethod(req, "")
	if handler != nil {
		return requestMapping.getHandlerExecutionChain(handler, req)
	}
	return nil
}

func (requestMapping RequestHandlerMapping) lookupHandlerMethod(req HttpRequest, lookupPath string) HandlerMethod {
	requestMatches := make([]RequestMatch, 0)
	directPathMappings, err := requestMapping.mappingRegistry.FindMappingsByUrl(lookupPath)
	if err == nil {
		requestMatches = append(requestMatches, requestMapping.getRequestMatches(req, directPathMappings)...)
	}
	/* todo complete this part which will match the given request with handlers */
	return nil
}

func (requestMapping RequestHandlerMapping) getRequestMatches(req HttpRequest, mappings []interface{}) []RequestMatch {
	matches := make([]RequestMatch, 0)
	for _, mapping := range mappings {
		match := mapping.(RequestMapping).MatchRequest(req)
		if match != nil {
			handlerMethod := requestMapping.mappingRegistry.GetMappings()[match]
			requestMatch := NewDefaultRequestMatch(mapping.(RequestMapping), NewSimpleHandlerMethod(handlerMethod))
			matches = append(matches, requestMatch)
		}
	}
	return matches
}

func (requestMapping RequestHandlerMapping) getHandlerExecutionChain(handlerMethod HandlerMethod, req HttpRequest) HandlerChain {
	return nil
}
