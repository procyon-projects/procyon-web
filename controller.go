package web

type Controller interface {
	RegisterHandlers(registry HandlerRegistry)
}

type MvcController interface {
	Routes() MvcRouterFunction
}

type MvcRequestHandler func(ctx MvcRequestContext)

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

type RestController interface {
	Routes() RestRouterFunction
}

type RestRequestHandler func(ctx RestRequestContext)

type RestRouterFunction interface {
	Get(pattern string, handler RestRequestHandler) RestRouterFunction
	Put(pattern string, handler RestRequestHandler) RestRouterFunction
	Post(pattern string, handler RestRequestHandler) RestRouterFunction
	Delete(pattern string, handler RestRequestHandler) RestRouterFunction
	Options(pattern string, handler RestRequestHandler) RestRouterFunction
	Head(pattern string, handler RestRequestHandler) RestRouterFunction
	Patch(pattern string, handler RestRequestHandler) RestRouterFunction
}

func RestRoute(prefix ...string) RestRouterFunction {
	return nil
}
