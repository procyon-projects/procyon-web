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
	// if an instance of http request exist, get from pool
	req := httpRequestPool.Get().(HttpRequest)
	req.request = request
	req.clearAttributes()
	// if an instance of http response exist, get from pool
	res := httpResponsePool.Get().(HttpResponse)
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
	// when you're done with the instances, put them into pool
	httpRequestPool.Put(req)
	httpResponsePool.Put(res)
}

func newWebServer(context WebApplicationContext) (Server, error) {
	return &DefaultWebServer{
		router: NewSimpleRouter(context),
	}, nil
}
