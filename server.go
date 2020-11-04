package web

import (
	"github.com/google/uuid"
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

type ProcyonWebServer struct {
	router     Router
	properties *configure.WebServerProperties
}

func (server *ProcyonWebServer) SetProperties(properties *configure.WebServerProperties) {
	server.properties = properties
}

func (server *ProcyonWebServer) Run() error {
	return http.ListenAndServe(":"+strconv.Itoa(server.GetPort()), server)
}

func (server *ProcyonWebServer) Stop() error {
	return nil
}

func (server *ProcyonWebServer) GetPort() int {
	var port int
	if server.properties == nil {
		port = DefaultWebServerPort
	} else {
		port = server.properties.Port
	}
	return port
}

func (server *ProcyonWebServer) ServeHTTP(response http.ResponseWriter, request *http.Request) {
	switch request.Method {
	case http.MethodGet:
		server.router.DoGet(response, request)
	case http.MethodPost:
		server.router.DoPost(response, request)
	case http.MethodPut:
		server.router.DoPut(response, request)
	case http.MethodDelete:
		server.router.DoDelete(response, request)
	case http.MethodPatch:
		server.router.DoPatch(response, request)
	}
}

func newProcyonWebServer(context WebApplicationContext) (Server, error) {
	server := &ProcyonWebServer{
		router: NewProcyonRouter(context.(ConfigurableWebApplicationContext)),
	}
	return server, nil
}

func NewProcyonWebServerForBenchmark(handlerRegistry SimpleHandlerRegistry) *ProcyonWebServer {
	appId, _ := uuid.NewUUID()
	contextId, _ := uuid.NewUUID()
	ctx := NewProcyonServerApplicationContext(appId, contextId)

	server := &ProcyonWebServer{
		router: newProcyonRouterForBenchmark(ctx.BaseWebApplicationContext, handlerRegistry),
	}
	return server
}
