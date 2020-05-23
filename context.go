package web

import (
	"github.com/Rollcomp/procyon-context"
)

// WebApplicationContext
type ApplicationContext interface {
	context.ApplicationContext
}

// ConfigurableWebApplicationContext
type ConfigurableApplicationContext interface {
	ApplicationContext
	context.ConfigurableContext
}

// GenericWebApplicationContext
type GenericApplicationContext struct {
	*context.GenericApplicationContext
}

func NewGenericApplicationContext(configurableContextAdapter context.ConfigurableContextAdapter) *GenericApplicationContext {
	return &GenericApplicationContext{
		context.NewGenericApplicationContext(configurableContextAdapter),
	}
}

// WebServeApplicationContext
type ServerApplicationContext interface {
	context.ApplicationContext
	GetWebServer() Server
}

type ConfigurableServerApplicationContext interface {
	ServerApplicationContext
	context.ConfigurableContext
}

// ---------------------------------------------------

// ProcyonServerApplicationContext
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
}

func (ctx *ProcyonServerApplicationContext) createWebServer() error {
	server, err := newWebServer()
	if err != nil {
		return err
	}
	ctx.server = server
	return nil
}
