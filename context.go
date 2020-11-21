package web

import (
	"github.com/codnect/goo"
	configure "github.com/procyon-projects/procyon-configure"
	"github.com/procyon-projects/procyon-context"
	core "github.com/procyon-projects/procyon-core"
	"github.com/valyala/fasthttp"
	"net/http"
	"strconv"
)

type ProcyonServerApplicationContext struct {
	*context.BaseApplicationContext
	server Server
}

func NewProcyonServerApplicationContext(appId context.ApplicationId, contextId context.ContextId) *ProcyonServerApplicationContext {
	ctx := &ProcyonServerApplicationContext{}
	applicationContext := context.NewBaseApplicationContext(appId, contextId, ctx)
	ctx.BaseApplicationContext = applicationContext
	return ctx
}

func (ctx *ProcyonServerApplicationContext) GetWebServer() Server {
	return ctx.server
}

func (ctx *ProcyonServerApplicationContext) Configure() {
	ctx.BaseApplicationContext.Configure()
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
	server, err := newProcyonWebServer(ctx.BaseApplicationContext)
	if err != nil {
		return err
	}
	ctx.server = server
	return nil
}

type PathVariable struct {
	Key   string
	Value string
}

type WebRequestContext struct {
	// context
	contextIdBuffer        [36]byte
	contextIdStr           string
	fastHttpRequestContext *fasthttp.RequestCtx
	// handler
	handlerChain  *HandlerChain
	handlerIndex  int
	inMainHandler bool
	completedFlow bool
	// path variables
	pathVariables     [20]string
	pathVariableCount int
	// response and error
	responseEntity ResponseEntity
	err            error
	// other
	valueMap map[string]interface{}
}

func newWebRequestContext() interface{} {
	return &WebRequestContext{
		handlerIndex: -1,
	}
}

func (ctx *WebRequestContext) prepare() {
	core.GenerateUUID(ctx.contextIdBuffer[:])
	ctx.contextIdStr = core.BytesToStr(ctx.contextIdBuffer[:])
	ctx.responseEntity.contentType = DefaultMediaType
}

func (ctx *WebRequestContext) reset() {
	ctx.handlerChain = nil
	ctx.handlerIndex = -1
	ctx.inMainHandler = true
	ctx.pathVariableCount = 0
	ctx.valueMap = nil
	ctx.responseEntity.status = http.StatusOK
	ctx.responseEntity.body = nil
	ctx.responseEntity.contentType = DefaultMediaType
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

func (ctx *WebRequestContext) GetContextId() context.ContextId {
	return context.ContextId(ctx.contextIdStr)
}

func (ctx *WebRequestContext) Get(key string) interface{} {
	return ctx.valueMap[key]
}

func (ctx *WebRequestContext) Put(key string, value interface{}) {
	ctx.valueMap[key] = value
}

func (ctx *WebRequestContext) addPathVariableValue(pathVariableName string) {
	ctx.pathVariables[ctx.pathVariableCount] = pathVariableName
	ctx.pathVariableCount++
}

func (ctx *WebRequestContext) GetPathVariable(name string) string {
	return ""
}

func (ctx *WebRequestContext) GetRequestParameter(name string) string {
	return ""
}

func (ctx *WebRequestContext) GetHeaderValue(key string) string {
	return ""
}

func (ctx *WebRequestContext) GetRequest() interface{} {
	return nil
}

func (ctx *WebRequestContext) SetStatus(status int) ResponseBodyBuilder {
	ctx.responseEntity.status = status
	return ctx
}

func (ctx *WebRequestContext) SetBody(body interface{}) ResponseBodyBuilder {
	ctx.responseEntity.body = body
	return ctx
}

func (ctx *WebRequestContext) SetContentType(mediaType MediaType) ResponseBodyBuilder {
	ctx.responseEntity.contentType = mediaType
	return ctx
}

func (ctx *WebRequestContext) AddHeader(key string, value string) ResponseHeaderBuilder {
	return ctx
}

func (ctx *WebRequestContext) GetStatus() int {
	return ctx.responseEntity.status
}

func (ctx *WebRequestContext) GetBody() interface{} {
	return ctx.responseEntity.body
}

func (ctx *WebRequestContext) GetContentType() MediaType {
	return ctx.responseEntity.contentType
}

func (ctx *WebRequestContext) Ok() ResponseBodyBuilder {
	ctx.responseEntity.status = http.StatusOK
	return ctx
}

func (ctx *WebRequestContext) NotFound() ResponseHeaderBuilder {
	ctx.responseEntity.status = http.StatusNotFound
	return ctx
}

func (ctx *WebRequestContext) NoContent() ResponseHeaderBuilder {
	ctx.responseEntity.status = http.StatusNoContent
	return ctx
}

func (ctx *WebRequestContext) BadRequest() ResponseBodyBuilder {
	ctx.responseEntity.status = http.StatusBadRequest
	return ctx
}

func (ctx *WebRequestContext) Accepted() ResponseBodyBuilder {
	ctx.responseEntity.status = http.StatusAccepted
	return ctx
}

func (ctx *WebRequestContext) Created(location string) ResponseBodyBuilder {
	ctx.responseEntity.status = http.StatusCreated
	return ctx
}

func (ctx *WebRequestContext) GetError() error {
	return nil
}

func (ctx *WebRequestContext) SetError(err error) {

}
