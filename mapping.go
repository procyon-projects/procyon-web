package web

type MappingRegistry interface {
	Register(path string, method RequestMethod, handlerChain *HandlerChain)
	Find(ctx *WebRequestContext)
}

type RequestMappingRegistry struct {
	routerTree *RouterTree
}

func NewRequestMappingRegistry() RequestMappingRegistry {
	return RequestMappingRegistry{
		routerTree: newRouterTree(),
	}
}

func (registry RequestMappingRegistry) Register(path string, method RequestMethod, handlerChain *HandlerChain) {
	registry.routerTree.AddRoute(path, method, handlerChain)
}

func (registry RequestMappingRegistry) Find(ctx *WebRequestContext) {
	registry.routerTree.Get(ctx)
}

type HandlerMapping interface {
	RegisterHandlerMethod(path string, method RequestMethod, handlerFunc RequestHandlerFunction)
	GetHandlerChain(ctx *WebRequestContext)
}

type RequestHandlerMapping struct {
	mappingRegistry     MappingRegistry
	interceptorRegistry HandlerInterceptorRegistry
}

func NewRequestHandlerMapping(mappingRegistry MappingRegistry, interceptorRegistry HandlerInterceptorRegistry) RequestHandlerMapping {
	return RequestHandlerMapping{
		mappingRegistry:     mappingRegistry,
		interceptorRegistry: interceptorRegistry,
	}
}

func (requestMapping RequestHandlerMapping) RegisterHandlerMethod(path string, method RequestMethod, handlerFunc RequestHandlerFunction) {
	requestMapping.mappingRegistry.Register(path, method, NewHandlerChain(handlerFunc, requestMapping.interceptorRegistry))
}

func (requestMapping RequestHandlerMapping) GetHandlerChain(ctx *WebRequestContext) {
	requestMapping.mappingRegistry.Find(ctx)
}
