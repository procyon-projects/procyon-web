package web

import (
	"errors"
	"github.com/codnect/goo"
	"github.com/procyon-projects/procyon-configure"
	"net/http"
	"strconv"
)

type Server interface {
	Run() error
	Stop() error
	GetPort() int
}

type DefaultWebServer struct {
	router     Router
	properties *configure.WebServerProperties
}

func (server *DefaultWebServer) setProperties(properties *configure.WebServerProperties) {
	server.properties = properties
}

func (server *DefaultWebServer) Run() error {
	port := server.properties.Port
	return http.ListenAndServe(":"+strconv.Itoa(port), server)
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
		_ = server.router.DoGet(res, req)
	case http.MethodPost:
		_ = server.router.DoPost(res, req)
	case http.MethodPut:
		_ = server.router.DoPut(res, req)
	case http.MethodDelete:
		_ = server.router.DoDelete(res, req)
	case http.MethodPatch:
		_ = server.router.DoPatch(res, req)
	}
	// when you're done with the instances, put them into pool
	httpRequestPool.Put(req)
	httpResponsePool.Put(res)
}

func newWebServer(context WebApplicationContext) (Server, error) {
	serverProperties := context.GetSharedPeaType(goo.GetType((*configure.WebServerProperties)(nil)))
	if serverProperties == nil {
		return nil, errors.New("an instance of configure.WebServerProperties not found")
	}
	server := &DefaultWebServer{
		router: NewSimpleRouter(context),
	}
	server.setProperties(serverProperties.(*configure.WebServerProperties))
	return server, nil
}
