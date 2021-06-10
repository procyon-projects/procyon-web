package web

import (
	"github.com/procyon-projects/goo"
	configure "github.com/procyon-projects/procyon-configure"
	"github.com/procyon-projects/procyon-context"
	core "github.com/procyon-projects/procyon-core"
	peas "github.com/procyon-projects/procyon-peas"
	"github.com/valyala/fasthttp"
	"net/http"
	"strconv"
	"strings"
	"unsafe"
)

type ServerApplicationContext struct {
	*context.BaseApplicationContext
	server Server
}

func NewProcyonServerApplicationContext(appId context.ApplicationId, contextId context.ContextId) *ServerApplicationContext {
	ctx := &ServerApplicationContext{}
	applicationContext := context.NewBaseApplicationContext(appId, contextId, ctx)
	ctx.BaseApplicationContext = applicationContext
	return ctx
}

func (ctx *ServerApplicationContext) GetWebServer() Server {
	return ctx.server
}

func (ctx *ServerApplicationContext) Configure() {
	ctx.BaseApplicationContext.Configure()
}

func (ctx *ServerApplicationContext) OnConfigure() {
	ctx.initializeInterceptors()
	_ = ctx.createWebServer()
}

func (ctx *ServerApplicationContext) initializeInterceptors() {
	peaFactory := ctx.BaseApplicationContext.GetPeaFactory()
	peaDefinitionRegistry := peaFactory.(peas.PeaDefinitionRegistry)
	peaNames := peaDefinitionRegistry.GetPeaDefinitionNames()

	for _, peaName := range peaNames {
		peaDefinition := peaDefinitionRegistry.GetPeaDefinition(peaName)
		if peaDefinition != nil && !ctx.isHandlerInterceptor(peaDefinition.GetPeaType()) {
			continue
		}
		peaFactory.GetPea(peaName)
	}
}

func (ctx *ServerApplicationContext) isHandlerInterceptor(typ goo.Type) bool {
	peaType := typ
	if peaType.IsFunction() {
		peaType = peaType.ToFunctionType().GetFunctionReturnTypes()[0]
	}

	if peaType.IsStruct() {
		structType := peaType.ToStructType()
		if structType.Implements(goo.GetType((*HandlerInterceptorBefore)(nil)).ToInterfaceType()) {
			return true
		} else if structType.Implements(goo.GetType((*HandlerInterceptorAfter)(nil)).ToInterfaceType()) {
			return true
		} else if structType.Implements(goo.GetType((*HandlerInterceptorAfterCompletion)(nil)).ToInterfaceType()) {
			return true
		}
	}
	return false
}

func (ctx *ServerApplicationContext) FinishConfigure() {
	logger := ctx.GetLogger()
	startedChannel := make(chan bool, 1)
	go func() {
		serverProperties := ctx.GetSharedPeaType(goo.GetType((*configure.WebServerProperties)(nil)))
		ctx.server.SetProperties(serverProperties.(*configure.WebServerProperties))
		logger.Info(ctx, "Procyon started on port(s): "+strconv.Itoa(int(ctx.GetWebServer().GetPort())))
		startedChannel <- true
		ctx.server.Run()
	}()
	<-startedChannel
}

func (ctx *ServerApplicationContext) createWebServer() error {
	ctx.server = newProcyonWebServer(ctx.BaseApplicationContext)
	return nil
}

type PathVariable struct {
	Key   string
	Value string
}

type WebRequestContext struct {
	router *ProcyonRouter
	// context
	contextIdBuffer        [36]byte
	contextIdStr           string
	fastHttpRequestContext *fasthttp.RequestCtx
	// cache
	path []byte
	args *fasthttp.Args
	uri  *fasthttp.URI
	// handler
	handlerChain *HandlerChain
	handlerIndex int
	// path variables
	pathVariables     [20]string
	pathVariableCount int
	// response and error
	responseWriter ResponseWriter
	responseEntity ResponseEntity
	httpError      *HTTPError
	internalError  error
	// other
	valueMap  map[string]interface{}
	canceled  bool
	completed bool
	crashed   bool
}

func (ctx *WebRequestContext) prepare(generateContextId bool) {
	if generateContextId {
		core.GenerateUUID(ctx.contextIdBuffer[:])
		ctx.contextIdStr = core.BytesToStr(ctx.contextIdBuffer[:])
	}
}

func (ctx *WebRequestContext) reset() {
	ctx.httpError = nil
	ctx.internalError = nil
	ctx.handlerChain = nil
	ctx.crashed = false
	ctx.canceled = false
	ctx.completed = false
	ctx.path = nil
	ctx.uri = nil
	ctx.args = nil
	ctx.handlerIndex = 0
	ctx.pathVariableCount = 0
	ctx.valueMap = nil
	ctx.responseEntity.status = http.StatusOK
	ctx.responseEntity.model = nil
	ctx.responseEntity.contentType = DefaultMediaType
	ctx.responseEntity.location = ""
}

func (ctx *WebRequestContext) writeResponse() {
	err := ctx.router.responseBodyWriter.WriteResponseBody(ctx, ctx.responseWriter)
	if err != nil {
		panic(err)
	}

	ctx.fastHttpRequestContext.SetStatusCode(ctx.responseEntity.status)

	if ctx.responseEntity.status == http.StatusCreated && ctx.responseEntity.location != "" {
		ctx.fastHttpRequestContext.Response.Header.Add(fasthttp.HeaderLocation, ctx.responseEntity.location)
	}

	switch ctx.responseEntity.contentType {
	case MediaTypeApplicationJson:
		ctx.fastHttpRequestContext.SetContentType(MediaTypeApplicationJsonValue)
	case MediaTypeApplicationTextHtml:
		ctx.fastHttpRequestContext.SetContentType(MediaTypeApplicationTextHtmlValue)
	default:
		ctx.fastHttpRequestContext.SetContentType(MediaTypeApplicationXmlValue)
	}
}

func (ctx *WebRequestContext) invoke() {
	if ctx.router.recoveryActive {
		defer ctx.router.errorHandlerManager.Recover(ctx)
		ctx.invokeHandlers()
	} else {
		ctx.invokeHandlers()
	}
}

func (ctx *WebRequestContext) invokeHandlers() {
next:
	if ctx.handlerIndex > ctx.handlerChain.handlerEndIndex {
		return
	}

	ctx.handlerChain.handlers[ctx.handlerIndex](ctx)
	if ctx.handlerIndex < ctx.handlerChain.handlerIndex && ctx.canceled {
		ctx.handlerIndex = ctx.handlerChain.afterCompletionStartIndex - 1
	}

	ctx.handlerIndex++
	if ctx.handlerIndex == ctx.handlerChain.afterCompletionStartIndex {
		if ctx.internalError == nil && ctx.httpError != nil {
			ctx.router.errorHandlerManager.JustHandleError(ctx.httpError, ctx)
		}
		ctx.writeResponse()
		ctx.completed = true
	}

	goto next
}

func (ctx *WebRequestContext) Cancel() {
	if ctx.handlerIndex < ctx.handlerChain.handlerIndex {
		ctx.canceled = true
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

func (ctx *WebRequestContext) GetRequestBody() []byte {
	return ctx.fastHttpRequestContext.Request.Body()
}

func (ctx *WebRequestContext) Validate(val interface{}) error {
	return ctx.router.validator.Validate(val)
}

func (ctx *WebRequestContext) BindRequest(request interface{}) error {
	return ctx.router.requestBinder.BindRequest(request, ctx)
}

func (ctx *WebRequestContext) GetResponseBody() []byte {
	return ctx.fastHttpRequestContext.Response.Body()
}

func (ctx *WebRequestContext) GetResponseHeader(key string) (string, bool) {
	val := ctx.fastHttpRequestContext.Response.Header.Peek(key)
	if val == nil {
		return "", false
	}
	return string(val), true
}

func (ctx *WebRequestContext) GetHTTPError() *HTTPError {
	return ctx.httpError
}

func (ctx *WebRequestContext) GetInternalError() error {
	return ctx.internalError
}

func (ctx *WebRequestContext) SetHTTPError(err *HTTPError) {
	if err != nil && ctx.handlerIndex <= ctx.handlerChain.handlerIndex {
		ctx.httpError = err
	}
}

func (ctx *WebRequestContext) IsSuccess() bool {
	return !ctx.crashed
}

func (ctx *WebRequestContext) IsCanceled() bool {
	return ctx.completed && ctx.canceled
}

func (ctx *WebRequestContext) IsCompleted() bool {
	ctx.fastHttpRequestContext.Response.SetBody()
	return ctx.completed && !ctx.canceled
}

type MvcRequestContext interface {
	GetRequestHeader(name string) (string, bool)
	GetQueryParameter(name string) (string, bool)
	GetPathVariable(name string) (string, bool)
	SetError(err error)
	AddObject(name string, value interface{})
	SetViewName(name string)
	GetViewName() string
	SetStatus(status int)
	GetStatus() int
}

type RestRequestContext interface {
	GetPath() string
	GetRequestHeader(name string) (string, bool)
	GetQueryParameter(name string) (string, bool)
	GetPathVariable(name string) (string, bool)
	Response() Response
}

type restRequestContext struct {
	// request context
	requestContext *fasthttp.RequestCtx
	args           *fasthttp.Args
	uri            *fasthttp.URI
	// path variables
	pathVariables     [20]string
	pathVariableCount int
	// handler
	handlerChain *HandlerChain
	handlerIndex int
	// response entity
	responseEntity *ResponseEntity
}

func (ctx *restRequestContext) clear() {
	ctx.args = nil
	ctx.uri = nil
	ctx.pathVariableCount = 0
	ctx.handlerIndex = 0
	ctx.handlerChain = nil
	ctx.responseEntity.reset()
}

func (ctx *restRequestContext) GetPath() string {
	if ctx.uri == nil {
		ctx.uri = ctx.requestContext.URI()
	}

	return *(*string)(unsafe.Pointer(&ctx.uri))
}

func (ctx *restRequestContext) GetRequestHeader(name string) (string, bool) {
	value := ctx.requestContext.Request.Header.Peek(name)

	if value == nil {
		return "", false
	}

	return *(*string)(unsafe.Pointer(&value)), true
}

func (ctx *restRequestContext) GetQueryParameter(name string) (string, bool) {
	if ctx.args == nil {
		ctx.args = ctx.requestContext.QueryArgs()
	}

	result := ctx.args.Peek(name)

	if result == nil {
		return "", false
	}

	return *(*string)(unsafe.Pointer(&result)), true
}

func (ctx *restRequestContext) GetPathVariable(name string) (string, bool) {
	for index, pathVariableName := range ctx.handlerChain.pathVariables {

		if pathVariableName == name {
			return ctx.pathVariables[index], true
		}

	}

	return "", false
}

func (ctx *restRequestContext) Response() Response {
	return ctx.responseEntity
}

func x(ctx RestRequestContext) error {
	ctx.Response().
		Status(http.StatusUnauthorized).
		Body("").
		ContentType(MediaTypeJson)

	return nil
}
