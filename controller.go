package web

type Controller interface {
	Initialize(registry HandlerInfoRegistry)
}
