package web

import (
	"net/http"
)

type Server interface {
	Run(args ...string) error
	Stop() error
	GetPort() int
}

type DefaultWebServer struct {
	router Router
}

func (server *DefaultWebServer) Run(args ...string) error {
	return http.ListenAndServe(":8080", server)
}

func (server *DefaultWebServer) Stop() error {
	return nil
}

func (server *DefaultWebServer) GetPort() int {
	return 8080
}

func (server *DefaultWebServer) ServeHTTP(response http.ResponseWriter, request *http.Request) {
	req := getHttpRequestFromPool()
	req.request = request
	res := getHttpResponseFromPool()
	res.responseWriter = response
	switch request.Method {
	case http.MethodGet:
		_ = server.router.DoDelete(res, req)
	case http.MethodPost:
		_ = server.router.DoPost(res, req)
	case http.MethodPut:
		_ = server.router.DoPut(res, req)
	case http.MethodDelete:
		_ = server.router.DoPut(res, req)
	case http.MethodPatch:
		_ = server.router.DoPatch(res, req)
	}
	putToPool(res, req)
}

func newWebServer(context ApplicationContext) (Server, error) {
	return &DefaultWebServer{
		router: NewSimpleRouter(context),
	}, nil
}
