package web

type Controller interface {
	RegisterHandlers(registry HandlerRegistry)
}

type MvcController interface {
	Routes() MvcRouterFunction
}

type MvcRouterFunction interface {
	Get(pattern string, handler func(ctx MvcRequestContext)) MvcRouterFunction
	Put(pattern string, handler func(ctx MvcRequestContext)) MvcRouterFunction
	Post(pattern string, handler func(ctx MvcRequestContext)) MvcRouterFunction
	Delete(pattern string, handler func(ctx MvcRequestContext)) MvcRouterFunction
	Options(pattern string, handler func(ctx MvcRequestContext)) MvcRouterFunction
	Head(pattern string, handler func(ctx MvcRequestContext)) MvcRouterFunction
	Patch(pattern string, handler func(ctx MvcRequestContext)) MvcRouterFunction
}

func MvcRoute(prefix ...string) MvcRouterFunction {
	return nil
}

type RestController interface {
	Routes() RestRouterFunction
}

type RestRouterFunction interface {
	Get(pattern string, handler func(ctx RestRequestContext)) RestRouterFunction
	Put(pattern string, handler func(ctx RestRequestContext)) RestRouterFunction
	Post(pattern string, handler func(ctx RestRequestContext)) RestRouterFunction
	Delete(pattern string, handler func(ctx RestRequestContext)) RestRouterFunction
	Options(pattern string, handler func(ctx RestRequestContext)) RestRouterFunction
	Head(pattern string, handler func(ctx RestRequestContext)) RestRouterFunction
	Patch(pattern string, handler func(ctx RestRequestContext)) RestRouterFunction
}

func RestRoute(prefix ...string) RestRouterFunction {
	return nil
}
