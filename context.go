package web

import (
	"github.com/codnect/goo"
	"github.com/google/uuid"
	"github.com/procyon-projects/procyon-configure"
	"github.com/procyon-projects/procyon-context"
	"net/http"
	"net/url"
	"strconv"
)

type PathVariable struct {
	Key   string
	Value string
}

type RequestContext interface {
	context.Context
	GetHeaderValue(key string) string
	GetPathVariables() []PathVariable
	GetPathVariable(name string) (string, bool)
	GetRequestParameter(name string) string
	GetRequestData() interface{}
}

type WebRequestContext struct {
	hasContextId      bool
	contextId         uuid.UUID
	request           *http.Request
	header            http.Header
	requestParamCache url.Values
	pathVariables     []PathVariable
	formCache         url.Values
	requestData       interface{}
}

func newWebRequestContext() interface{} {
	return &WebRequestContext{}
}

func (context *WebRequestContext) reset() {
	if !context.hasContextId {
		context.contextId, _ = uuid.NewUUID()
	}
	context.request = nil
	context.header = nil
	context.pathVariables = context.pathVariables[0:0]
	context.requestParamCache = nil
	context.formCache = nil
	context.requestData = nil
}

func (context *WebRequestContext) GetContextId() uuid.UUID {
	return context.contextId
}

func (context *WebRequestContext) GetHeaderValue(key string) string {
	if context.header == nil {
		context.header = context.request.Header
	}
	return context.header.Get(key)
}

func (context *WebRequestContext) GetPathVariables() []PathVariable {
	return context.pathVariables
}

func (context *WebRequestContext) GetPathVariable(name string) (string, bool) {
	for _, variable := range context.pathVariables {
		if variable.Key == name {
			return variable.Value, true
		}
	}
	return "", false
}

func (context *WebRequestContext) GetRequestParameter(name string) string {
	if context.requestParamCache == nil {
		context.requestParamCache = context.request.URL.Query()
	}
	return context.requestParamCache.Get(name)
}

func (context *WebRequestContext) GetRequestData() interface{} {
	return context.requestData
}

type WebApplicationContext interface {
	context.ApplicationContext
}

type ConfigurableWebApplicationContext interface {
	WebApplicationContext
	context.ConfigurableContext
}

type BaseWebApplicationContext struct {
	*context.BaseApplicationContext
}

func NewBaseWebApplicationContext(appId uuid.UUID, contextId uuid.UUID, configurableContextAdapter context.ConfigurableContextAdapter) *BaseWebApplicationContext {
	return &BaseWebApplicationContext{
		context.NewBaseApplicationContext(appId, contextId, configurableContextAdapter),
	}
}

type ServerApplicationContext interface {
	context.ApplicationContext
	GetWebServer() Server
}

type ConfigurableServerApplicationContext interface {
	ServerApplicationContext
	context.ConfigurableContext
}

type ProcyonServerApplicationContext struct {
	*BaseWebApplicationContext
	server Server
}

func NewProcyonServerApplicationContext(appId uuid.UUID, contextId uuid.UUID) *ProcyonServerApplicationContext {
	ctx := &ProcyonServerApplicationContext{}
	genericCtx := NewBaseWebApplicationContext(appId, contextId, ctx)
	ctx.BaseWebApplicationContext = genericCtx
	return ctx
}

func (ctx *ProcyonServerApplicationContext) GetWebServer() Server {
	return ctx.server
}

func (ctx *ProcyonServerApplicationContext) Configure() {
	ctx.BaseWebApplicationContext.BaseApplicationContext.Configure()
}

func (ctx *ProcyonServerApplicationContext) OnConfigure() {
	_ = ctx.createWebServer()
}

func (ctx *ProcyonServerApplicationContext) FinishConfigure() {
	logger := ctx.GetLogger()
	startedChannel := make(chan bool, 1)
	go func() {
		serverProperties := ctx.GetSharedPeaType(goo.GetType((*configure.WebServerProperties)(nil)))
		ctx.server.SetProperties(serverProperties.(*configure.WebServerProperties))
		logger.Info(ctx, "Procyon started on port(s): "+strconv.Itoa(ctx.GetWebServer().GetPort()))
		startedChannel <- true
		ctx.server.Run()
	}()
	<-startedChannel
}

func (ctx *ProcyonServerApplicationContext) createWebServer() error {
	server, err := newProcyonWebServer(ctx.BaseWebApplicationContext)
	if err != nil {
		return err
	}
	ctx.server = server
	return nil
}

func (ctx *ProcyonServerApplicationContext) cloneApplicationContext() {

}
