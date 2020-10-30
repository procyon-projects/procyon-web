package web

import (
	"errors"
	"github.com/codnect/goo"
	"sync"
)

const HandlerMappingUriVariableAttribute = "github.com.procyon.projects.procyon.handlermapping.urivariables"

type RequestMapping struct {
	methodRequestMatcher  MethodRequestMatcher
	paramsRequestMatcher  ParametersRequestMatcher
	patternRequestMatcher PatternRequestMatcher
}

func NewRequestMapping(methodRequestMatcher MethodRequestMatcher,
	paramsRequestMatcher ParametersRequestMatcher,
	patternRequestMatcher PatternRequestMatcher) *RequestMapping {
	return &RequestMapping{
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

func (registry RequestMappingRegistry) FindMappingsByUrl(path string) ([]interface{}, error) {
	if mappings, ok := registry.mappingUrlLookup[path]; ok {
		return mappings, nil
	}
	return nil, errors.New("not found matching")
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
	GetHandlerChain(req HttpRequest) HandlerChain
}

type RequestHandlerMapping struct {
	pathMatcher     PathMatcher
	mappingRegistry MappingRegistry
	interceptors    []HandlerInterceptor
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

func (requestMapping RequestHandlerMapping) GetHandlerChain(req HttpRequest) HandlerChain {
	handler := requestMapping.lookupHandlerMethod(req)
	if handler != nil {
		return requestMapping.getHandlerExecutionChain(handler, req)
	}
	return nil
}

func (requestMapping RequestHandlerMapping) lookupHandlerMethod(req HttpRequest) HandlerMethod {
	requestMatches := make([]RequestMatch, 0)
	lookupPath := req.GetPath()
	directPathMappings, err := requestMapping.mappingRegistry.FindMappingsByUrl(lookupPath)
	if err == nil {
		requestMatches = append(requestMatches, requestMapping.getRequestMatches(req, directPathMappings)...)
	}
	if len(requestMatches) == 0 {
		mappings := GetMapKeys(requestMapping.mappingRegistry.GetMappings())
		requestMatches = append(requestMatches, requestMapping.getRequestMatches(req, mappings)...)
	}
	if len(requestMatches) != 0 {
		match := requestMatches[0]
		mapping := match.GetMapping().(*RequestMapping)
		variables := requestMapping.pathMatcher.GetUriVariables(lookupPath, mapping.getPatternRequestMatcher().patterns[0])
		req.AddAttribute(HandlerMappingUriVariableAttribute, variables)
		return match.GetHandlerMethod()
	}
	return nil
}

func GetMapKeys(mapObj interface{}) []interface{} {
	argMapKeys := goo.GetType(mapObj).GetGoValue().MapKeys()
	mapKeys := make([]interface{}, len(argMapKeys))
	for i := 0; i < len(argMapKeys); i++ {
		mapKeys[i] = argMapKeys[i].Interface()
	}
	return mapKeys
}

func (requestMapping RequestHandlerMapping) getRequestMatches(req HttpRequest, mappings []interface{}) []RequestMatch {
	requestMatches := make([]RequestMatch, 0)
	for _, mapping := range mappings {
		match := mapping.(*RequestMapping).MatchRequest(req)
		if match != nil {
			handlerMethod := requestMapping.mappingRegistry.GetMappings()[mapping]
			requestMatch := NewDefaultRequestMatch(mapping.(*RequestMapping), NewSimpleHandlerMethod(handlerMethod))
			requestMatches = append(requestMatches, requestMatch)
		}
	}
	return requestMatches
}

func (requestMapping RequestHandlerMapping) getHandlerExecutionChain(handlerMethod HandlerMethod, req HttpRequest) HandlerChain {
	return NewHandlerExecutionChain(handlerMethod, WithInterceptors(requestMapping.interceptors))
}
