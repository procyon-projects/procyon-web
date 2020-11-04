package web

import (
	"github.com/codnect/goo"
	"github.com/google/uuid"
	"github.com/procyon-projects/procyon-configure"
	"github.com/procyon-projects/procyon-context"
	"strconv"
)

type PathVariable struct {
	Key   string
	Value string
}

type RequestContext interface {
	context.Context
	AddAttribute(key string, val interface{})
	HasAttribute(key string) bool
	GetAttribute(key string) interface{}
	GetPathVariables() []PathVariable
}

type WebRequestContext struct {
	contextId     uuid.UUID
	attributes    map[string]interface{}
	pathVariables []PathVariable
}

func newWebRequestContext() interface{} {
	return &WebRequestContext{
		attributes:    make(map[string]interface{}),
		pathVariables: make([]PathVariable, 0),
	}
}

func (context WebRequestContext) GetAppId() uuid.UUID {
	return uuid.UUID{}
}

func (context WebRequestContext) GetContextId() uuid.UUID {
	return context.contextId
}

func (context WebRequestContext) GetApplicationName() string {
	return ""
}

func (context WebRequestContext) GetStartupTimestamp() int64 {
	return 0
}

func (context WebRequestContext) AddAttribute(key string, val interface{}) {
	context.attributes[key] = val
}

func (context WebRequestContext) HasAttribute(key string) bool {
	if _, ok := context.attributes[key]; ok {
		return true
	}
	return false
}

func (context WebRequestContext) GetAttribute(key string) interface{} {
	if value, ok := context.attributes[key]; ok {
		return value
	}
	return nil
}

func (context WebRequestContext) GetPathVariables() []PathVariable {
	return context.pathVariables
}

func (context WebRequestContext) clear() {
	for key, _ := range context.attributes {
		delete(context.attributes, key)
	}
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
