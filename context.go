package web

import (
	"github.com/google/uuid"
	"github.com/procyon-projects/procyon-context"
)

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
	ctx.server.Run()
}

func (ctx *ProcyonServerApplicationContext) createWebServer() error {
	server, err := newWebServer(ctx.BaseWebApplicationContext)
	if err != nil {
		return err
	}
	ctx.server = server
	return nil
}

func (ctx *ProcyonServerApplicationContext) cloneApplicationContext() {

}
