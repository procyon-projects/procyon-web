package web

type Controller interface {
	RegisterHandlers(registry HandlerRegistry)
}

type MvcController interface {
	Routes() MvcRouterFunction
}

type MvcRequestHandler func(ctx MvcRequestContext) error

type MvcRouterFunction interface {
	Get(pattern string, handler MvcRequestHandler) MvcRouterFunction
	Put(pattern string, handler MvcRequestHandler) MvcRouterFunction
	Post(pattern string, handler MvcRequestHandler) MvcRouterFunction
	Delete(pattern string, handler MvcRequestHandler) MvcRouterFunction
	Options(pattern string, handler MvcRequestHandler) MvcRouterFunction
	Head(pattern string, handler MvcRequestHandler) MvcRouterFunction
	Patch(pattern string, handler MvcRequestHandler) MvcRouterFunction
}

func MvcRoute(prefix ...string) MvcRouterFunction {
	return nil
}

type RestRequestHandler func(ctx RestRequestContext) error

type RestRouterFunction interface {
	Get(pattern string, handler RestRequestHandler) RestRouterFunction
	Put(pattern string, handler RestRequestHandler) RestRouterFunction
	Post(pattern string, handler RestRequestHandler) RestRouterFunction
	Delete(pattern string, handler RestRequestHandler) RestRouterFunction
	Options(pattern string, handler RestRequestHandler) RestRouterFunction
	Head(pattern string, handler RestRequestHandler) RestRouterFunction
	Patch(pattern string, handler RestRequestHandler) RestRouterFunction
}

type RestController interface {
	Routes() RestRouterFunction
}

type restRouterFunction struct {
	method  HttpMethod
	pattern string
	handler RestRequestHandler
}

func newRestRouterFunction(method HttpMethod, pattern string, handler RestRequestHandler) restRouterFunction {
	if handler == nil {
		panic("handler cannot be nil")
	}

	return restRouterFunction{
		method:  method,
		pattern: pattern,
		handler: handler,
	}
}

type restRouterFunctions struct {
	prefix          string
	routerFunctions []restRouterFunction
	lastIndex       int
}

func RestRoute(prefix ...string) RestRouterFunction {
	routerFunctions := &restRouterFunctions{
		routerFunctions: make([]restRouterFunction, 0),
	}

	if len(prefix) > 0 {
		routerFunctions.prefix = prefix[0]
	}

	return routerFunctions
}

func (routes *restRouterFunctions) Get(pattern string, handler RestRequestHandler) RestRouterFunction {
	routes.routerFunctions = append(routes.routerFunctions, newRestRouterFunction(HttpMethodGet, pattern, handler))
	return routes
}

func (routes *restRouterFunctions) Put(pattern string, handler RestRequestHandler) RestRouterFunction {
	routes.routerFunctions = append(routes.routerFunctions, newRestRouterFunction(HttpMethodPut, pattern, handler))
	return routes
}

func (routes *restRouterFunctions) Post(pattern string, handler RestRequestHandler) RestRouterFunction {
	routes.routerFunctions = append(routes.routerFunctions, newRestRouterFunction(HttpMethodPost, pattern, handler))
	return routes
}

func (routes *restRouterFunctions) Delete(pattern string, handler RestRequestHandler) RestRouterFunction {
	routes.routerFunctions = append(routes.routerFunctions, newRestRouterFunction(HttpMethodDelete, pattern, handler))
	return routes
}

func (routes *restRouterFunctions) Options(pattern string, handler RestRequestHandler) RestRouterFunction {
	routes.routerFunctions = append(routes.routerFunctions, newRestRouterFunction(HttpMethodOptions, pattern, handler))
	return routes
}

func (routes *restRouterFunctions) Head(pattern string, handler RestRequestHandler) RestRouterFunction {
	routes.routerFunctions = append(routes.routerFunctions, newRestRouterFunction(HttpMethodHead, pattern, handler))
	return routes
}

func (routes *restRouterFunctions) Patch(pattern string, handler RestRequestHandler) RestRouterFunction {
	routes.routerFunctions = append(routes.routerFunctions, newRestRouterFunction(HttpMethodPatch, pattern, handler))
	return routes
}
