package web

import "sync"

type HandlerMethodParameterResolver interface {
	SupportsParameter(parameter HandlerMethodParameter) bool
	ResolveParameter(parameter HandlerMethodParameter, request HttpRequest) (interface{}, error)
}

type HandlerMethodParameterResolvers struct {
	parameterResolverCache map[int]HandlerMethodParameterResolver
	resolvers              []HandlerMethodParameterResolver
	cacheMutex             sync.Mutex
}

func NewHandlerMethodParameterResolvers() *HandlerMethodParameterResolvers {
	return &HandlerMethodParameterResolvers{
		parameterResolverCache: make(map[int]HandlerMethodParameterResolver),
		resolvers:              make([]HandlerMethodParameterResolver, 0),
	}
}

func (parameterResolvers *HandlerMethodParameterResolvers) SupportsParameter(parameter HandlerMethodParameter) bool {
	var cacheResolver HandlerMethodParameterResolver
	parameterResolvers.cacheMutex.Lock()
	cacheResolver = parameterResolvers.parameterResolverCache[parameter.HashCode()]
	parameterResolvers.cacheMutex.Unlock()
	if cacheResolver != nil {
		return true
	}
	resolvers := parameterResolvers.resolvers
	for _, resolver := range resolvers {
		if resolver.SupportsParameter(parameter) {
			parameterResolvers.cacheMutex.Lock()
			parameterResolvers.parameterResolverCache[parameter.HashCode()] = cacheResolver
			parameterResolvers.cacheMutex.Unlock()
			return true
		}
	}
	return false
}

func (parameterResolvers *HandlerMethodParameterResolvers) ResolveParameter(parameter HandlerMethodParameter, request HttpRequest) (interface{}, error) {
	resolver := parameterResolvers.findParameterResolver(parameter)
	if resolver == nil {
		return nil, NewNoHandlerParameterResolverError("Parameter resolver not found")
	}
	return resolver.ResolveParameter(parameter, request)
}

func (parameterResolvers *HandlerMethodParameterResolvers) findParameterResolver(parameter HandlerMethodParameter) HandlerMethodParameterResolver {
	var cacheResolver HandlerMethodParameterResolver
	parameterResolvers.cacheMutex.Lock()
	cacheResolver = parameterResolvers.parameterResolverCache[parameter.HashCode()]
	parameterResolvers.cacheMutex.Unlock()
	if cacheResolver != nil {
		return cacheResolver
	}
	resolvers := parameterResolvers.resolvers
	for _, resolver := range resolvers {
		if resolver.SupportsParameter(parameter) {
			parameterResolvers.cacheMutex.Lock()
			parameterResolvers.parameterResolverCache[parameter.HashCode()] = cacheResolver
			parameterResolvers.cacheMutex.Unlock()
			return resolver
		}
	}
	return nil
}

func (parameterResolvers *HandlerMethodParameterResolvers) AddMethodParameterResolver(resolvers ...HandlerMethodParameterResolver) {
	parameterResolvers.resolvers = append(parameterResolvers.resolvers, resolvers...)
}

type RequestMethodParameterResolver struct {
}

func NewRequestMethodParameterResolver() RequestMethodParameterResolver {
	return RequestMethodParameterResolver{}
}

func (r RequestMethodParameterResolver) SupportsParameter(parameter HandlerMethodParameter) bool {
	return true
}

func (r RequestMethodParameterResolver) ResolveParameter(parameter HandlerMethodParameter, request HttpRequest) (interface{}, error) {
	return nil, nil
}
