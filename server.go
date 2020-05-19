package web

import (
	"net/http"
)

type Server interface {
	Run(args ...string) error
}

type DefaultWebServer struct {
	handler Handler
}

func (server *DefaultWebServer) Run(args ...string) error {
	return http.ListenAndServe(":8080", server)
}

func (server *DefaultWebServer) ServeHTTP(res http.ResponseWriter, request *http.Request) {
	switch request.Method {
	case http.MethodGet:
		server.handler.DoDelete(res, request)
	case http.MethodPost:
		server.handler.DoPost(res, request)
	case http.MethodPut:
		server.handler.DoPut(res, request)
	case http.MethodDelete:
		server.handler.DoPut(res, request)
	case http.MethodPatch:
		server.handler.DoPatch(res, request)
	}
	res.WriteHeader(200)
}

func GetWebServer() (Server, error) {
	return &DefaultWebServer{
		handler: NewDefaultHandler(),
	}, nil
}
