package web

import (
	"errors"
	"net/http"
	"sync"
)

const HandlerMappingLookupPath = "github.com.procyon.projects.procyon.HandlerMapping.LookupPath"

type RequestMapping struct {
	methodRequestMatcher  MethodRequestMatcher
	patternRequestMatcher PatternRequestMatcher
}

func NewRequestMapping(methodRequestMatcher MethodRequestMatcher,
	patternRequestMatcher PatternRequestMatcher) *RequestMapping {
	return &RequestMapping{
		methodRequestMatcher,
		patternRequestMatcher,
	}
}

func (mappingInfo RequestMapping) getMethodRequestMatcher() MethodRequestMatcher {
	return mappingInfo.methodRequestMatcher
}

func (mappingInfo RequestMapping) getPatternRequestMatcher() PatternRequestMatcher {
	return mappingInfo.patternRequestMatcher
}

func (mappingInfo RequestMapping) MatchRequest(requestContext RequestContext, req *http.Request) bool {
	method := mappingInfo.methodRequestMatcher.MatchRequest(requestContext, req)
	if !method {
		return false
	} else {
		pattern := mappingInfo.patternRequestMatcher.MatchRequest(requestContext, req)
		if !pattern {
			return false
		}
	}
	return true
}

func (mappingInfo RequestMapping) hashCode() int {
	return 31*mappingInfo.patternRequestMatcher.hashCode() +
		mappingInfo.methodRequestMatcher.hashCode()
}

type MappingRegistry interface {
	Register(handlerName string, mapping interface{}, fun RequestHandlerFunc) error
	GetMappings() map[interface{}]HandlerMethod
	FindMappingsByUrl(path string) ([]interface{}, bool)
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
	requestMapping := mapping.(*RequestMapping)
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

func (registry RequestMappingRegistry) FindMappingsByUrl(path string) ([]interface{}, bool) {
	if mappings, ok := registry.mappingUrlLookup[path]; ok {
		return mappings, true
	}
	return nil, false
}

func (registry RequestMappingRegistry) findPurePaths(mapping *RequestMapping) {
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
	GetHandler(requestContext RequestContext, req *http.Request) HandlerMethod
}

type RequestHandlerMapping struct {
	pathMatcher     PathMatcher
	mappingRegistry MappingRegistry
	mu              sync.Mutex
}

func NewRequestHandlerMapping(pathMatcher PathMatcher, mappingRegistry MappingRegistry) RequestHandlerMapping {
	return RequestHandlerMapping{
		pathMatcher:     pathMatcher,
		mappingRegistry: mappingRegistry,
	}
}

func (requestMapping RequestHandlerMapping) RegisterHandlerMethod(handlerName string, mapping interface{}, fun RequestHandlerFunc) {
	requestMapping.mappingRegistry.Register(handlerName, mapping, fun)
}

func (requestMapping RequestHandlerMapping) GetHandler(requestContext RequestContext, req *http.Request) HandlerMethod {
	return requestMapping.lookupHandlerMethod(requestContext, req)
}

func (requestMapping RequestHandlerMapping) lookupHandlerMethod(requestContext RequestContext, req *http.Request) HandlerMethod {
	directPathMappings, ok := requestMapping.mappingRegistry.FindMappingsByUrl(req.URL.Path)
	if ok {
		_, handlerMethod := requestMapping.getRequestMatches(requestContext, req, directPathMappings)
		return handlerMethod
	}
	requestMatch, handlerMethod := requestMapping.scanRequestMatches(requestContext, req)
	if requestMatch != nil {
		return handlerMethod
	}
	return nil
}

func (requestMapping RequestHandlerMapping) getRequestMatches(requestContext RequestContext, req *http.Request, mappings []interface{}) (*RequestMapping, HandlerMethod) {
	for _, mapping := range mappings {
		if mapping.(*RequestMapping).MatchRequest(requestContext, req) {
			handlerMethod := requestMapping.mappingRegistry.GetMappings()[mapping]
			return mapping.(*RequestMapping), handlerMethod
		}
	}
	return nil, nil
}

func (requestMapping RequestHandlerMapping) scanRequestMatches(requestContext RequestContext, req *http.Request) (*RequestMapping, HandlerMethod) {
	for mapping, handlerMethod := range requestMapping.mappingRegistry.GetMappings() {
		if mapping.(*RequestMapping).MatchRequest(requestContext, req) {
			return mapping.(*RequestMapping), handlerMethod
		}
	}
	return nil, nil
}
