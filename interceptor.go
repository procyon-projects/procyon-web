package web

import "net/http"

type HandlerInterceptor interface {
	HandleBefore(handler interface{}, requestContext RequestContext, res http.ResponseWriter, req *http.Request)
	HandleAfter(handler interface{}, requestContext RequestContext, res http.ResponseWriter, req *http.Request)
}
