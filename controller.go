package web

type Controller interface {
	RegisterHandlers(registry HandlerInfoRegistry)
}
