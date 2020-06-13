package web

type HandlerMethodParameterResolver interface {
	SupportsParameter(parameter HandlerMethodParameter) bool
	ResolveParameter(parameter HandlerMethodParameter, request HttpRequest) (interface{}, error)
}

type HandlerMethodParameterResolvers struct {
	resolvers []HandlerMethodParameterResolver
}

func NewHandlerMethodParameterResolvers() *HandlerMethodParameterResolvers {
	return &HandlerMethodParameterResolvers{
		make([]HandlerMethodParameterResolver, 0),
	}
}

func (r *HandlerMethodParameterResolvers) SupportsParameter(parameter HandlerMethodParameter) bool {
	resolvers := r.resolvers
	for _, resolver := range resolvers {
		if resolver.SupportsParameter(parameter) {
			return true
		}
	}
	return false
}

func (r *HandlerMethodParameterResolvers) ResolveParameter(parameter HandlerMethodParameter, request HttpRequest) (interface{}, error) {
	resolver := r.findParameterResolver(parameter)
	if resolver == nil {
		return nil, NewNoHandlerParameterResolverError("Parameter resolver not found")
	}
	return resolver.ResolveParameter(parameter, request)
}

func (r *HandlerMethodParameterResolvers) findParameterResolver(parameter HandlerMethodParameter) HandlerMethodParameterResolver {
	resolvers := r.resolvers
	for _, resolver := range resolvers {
		if resolver.SupportsParameter(parameter) {
			return resolver
		}
	}
	return nil
}

func (r *HandlerMethodParameterResolvers) AddMethodParameterResolver(resolvers ...HandlerMethodParameterResolver) {
	r.resolvers = append(r.resolvers, resolvers...)
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
