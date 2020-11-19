package web

import (
	"github.com/codnect/goo"
	"github.com/procyon-projects/procyon-configure"
	"github.com/procyon-projects/procyon-context"
	core "github.com/procyon-projects/procyon-core"
	"github.com/valyala/fasthttp"
	"strconv"
)

type PathVariable struct {
	Key   string
	Value string
}

type WebRequestContext struct {
	contextIdBuffer        [36]byte
	contextStr             string
	handlerChain           *HandlerChain
	fastHttpRequestContext *fasthttp.RequestCtx
	pathVariables          [20]string
	paramCount             int
	responseEntity         *ResponseEntity
	handlerIndex           int
	inMainHandler          bool
	completedFlow          bool
	err                    error
	needReset              bool
}

func newWebRequestContext() interface{} {
	return &WebRequestContext{
		handlerIndex: -1,
	}
}

func (ctx *WebRequestContext) GetContextId() context.ContextId {
	return "ctx.contextStr"
}

func (ctx *WebRequestContext) reset() {
	ctx.handlerChain = nil
	ctx.handlerIndex = -1
	ctx.inMainHandler = true
	ctx.paramCount = 0
}

func (ctx *WebRequestContext) prepare() {
	core.GenerateUUID(ctx.contextIdBuffer[:])
	ctx.contextStr = bytesToStr(ctx.contextIdBuffer[:])
}

func (ctx *WebRequestContext) Next() {
	if ctx.inMainHandler {
		return
	}
	if ctx.handlerIndex >= ctx.handlerChain.handlerEndIndex {
		return
	}
	ctx.handlerIndex++
	if ctx.handlerIndex < ctx.handlerChain.handlerIndex {
		ctx.handlerChain.allHandlers[ctx.handlerIndex](ctx)
		return
	} else if ctx.handlerIndex == ctx.handlerChain.handlerIndex {
		ctx.inMainHandler = true
		ctx.handlerChain.allHandlers[ctx.handlerIndex](ctx)
		ctx.handlerIndex++
		ctx.inMainHandler = false
	}
	if ctx.handlerIndex < ctx.handlerChain.afterStartIndex {
		ctx.handlerChain.allHandlers[ctx.handlerIndex](ctx)
		return
	} else if ctx.handlerIndex < ctx.handlerChain.afterCompletionStartIndex {
		ctx.completedFlow = true
		ctx.handlerChain.allHandlers[ctx.handlerIndex](ctx)
	}
}

func (ctx *WebRequestContext) addPathVariableValue(pathVariableName string) {
	ctx.pathVariables[ctx.paramCount] = pathVariableName
	ctx.paramCount++
}

func (ctx *WebRequestContext) GetHeaderValue(key string) string {
	return ""
}

func (ctx *WebRequestContext) GetPathVariables() []PathVariable {
	return nil
}

func (ctx *WebRequestContext) GetPathVariable(name string) (string, bool) {
	/*	for _, variable := range context.pathVariables {
		if variable.Key == name {
			return variable.Value, true
		}
	}*/
	return "", false
}

func (ctx *WebRequestContext) GetRequestParameter(name string) string {
	return ""
}

func (ctx *WebRequestContext) GetRequestData() interface{} {
	return nil
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

func NewBaseWebApplicationContext(appId context.ApplicationId, contextId context.ContextId, configurableContextAdapter context.ConfigurableContextAdapter) *BaseWebApplicationContext {
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

func NewProcyonServerApplicationContext(appId context.ApplicationId, contextId context.ContextId) *ProcyonServerApplicationContext {
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
