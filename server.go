package web

import (
	"github.com/google/uuid"
	"github.com/procyon-projects/procyon-configure"
	"github.com/procyon-projects/procyon-context"
	"github.com/valyala/fasthttp"
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
	return fasthttp.ListenAndServe(":"+strconv.Itoa(server.GetPort()), server.Handle)
}

func (server *ProcyonWebServer) Handle(ctx *fasthttp.RequestCtx) {
	server.router.Route(ctx)
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

func newProcyonWebServer(ctx context.ApplicationContext) Server {
	server := &ProcyonWebServer{
		router: NewProcyonRouter(ctx.(context.ConfigurableApplicationContext)),
	}
	return server
}

func NewProcyonWebServerForBenchmark(handlerRegistry SimpleHandlerRegistry) *ProcyonWebServer {
	appId := uuid.New()
	contextId := uuid.New()
	ctx := NewProcyonServerApplicationContext(context.ApplicationId(appId.String()), context.ContextId(contextId.String()))

	server := &ProcyonWebServer{
		router: newProcyonRouterForBenchmark(ctx.BaseApplicationContext, handlerRegistry),
	}
	return server
}
