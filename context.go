package web

import (
	"github.com/procyon-projects/procyon-context"
)

type ApplicationContext interface {
	context.ApplicationContext
}

type ConfigurableApplicationContext interface {
	ApplicationContext
	context.ConfigurableContext
}

type GenericApplicationContext struct {
	*context.GenericApplicationContext
}

func NewGenericApplicationContext(configurableContextAdapter context.ConfigurableContextAdapter) *GenericApplicationContext {
	return &GenericApplicationContext{
		context.NewGenericApplicationContext(configurableContextAdapter),
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
	*GenericApplicationContext
	server Server
}

func NewProcyonServerApplicationContext() *ProcyonServerApplicationContext {
	ctx := &ProcyonServerApplicationContext{}
	genericCtx := NewGenericApplicationContext(ctx)
	ctx.GenericApplicationContext = genericCtx
	return ctx
}

func (ctx *ProcyonServerApplicationContext) GetWebServer() Server {
	return ctx.server
}

func (ctx *ProcyonServerApplicationContext) Configure() {
	ctx.GenericApplicationContext.GenericApplicationContext.Configure()
}

func (ctx *ProcyonServerApplicationContext) OnConfigure() {
	_ = ctx.createWebServer()
	ctx.server.Run()
}

func (ctx *ProcyonServerApplicationContext) createWebServer() error {
	server, err := newWebServer(ctx.GenericApplicationContext)
	if err != nil {
		return err
	}
	ctx.server = server
	return nil
}
