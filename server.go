package web

import (
	"github.com/procyon-projects/procyon-configure"
	"net/http"
	"strconv"
)

type Server interface {
	Run() error
	Stop() error
	SetProperties(properties *configure.WebServerProperties)
	GetPort() int
}

const DefaultWebServerPort = 8080

type DefaultWebServer struct {
	router     Router
	properties *configure.WebServerProperties
}

func (server *DefaultWebServer) SetProperties(properties *configure.WebServerProperties) {
	server.properties = properties
}

func (server *DefaultWebServer) Run() error {
	return http.ListenAndServe(":"+strconv.Itoa(server.GetPort()), server)
}

func (server *DefaultWebServer) Stop() error {
	return nil
}

func (server *DefaultWebServer) GetPort() int {
	var port int
	if server.properties == nil {
		port = DefaultWebServerPort
	} else {
		port = server.properties.Port
	}
	return port
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
	server := &DefaultWebServer{
		router: NewSimpleRouter(context),
	}
	return server, nil
}
