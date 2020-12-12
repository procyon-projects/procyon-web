package web

import (
	"encoding/xml"
	"github.com/codnect/goo"
	json "github.com/json-iterator/go"
	configure "github.com/procyon-projects/procyon-configure"
	"github.com/procyon-projects/procyon-context"
	core "github.com/procyon-projects/procyon-core"
	"github.com/valyala/fasthttp"
	"net/http"
	"reflect"
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
	ctx.server = newProcyonWebServer(ctx.BaseApplicationContext)
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
	// cache
	path []byte
	args *fasthttp.Args
	uri  *fasthttp.URI
	// handler
	handlerChain  *HandlerChain
	handlerIndex  int
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
		handlerIndex: 0,
		valueMap:     make(map[string]interface{}),
	}
}

func (ctx *WebRequestContext) prepare(generateContextId bool) {
	if generateContextId {
		core.GenerateUUID(ctx.contextIdBuffer[:])
		ctx.contextIdStr = core.BytesToStr(ctx.contextIdBuffer[:])
	}
}

func (ctx *WebRequestContext) reset() {
	ctx.path = nil
	ctx.uri = nil
	ctx.args = nil
	ctx.handlerIndex = 0
	ctx.pathVariableCount = 0
	ctx.valueMap = nil
	ctx.responseEntity.status = http.StatusOK
	ctx.responseEntity.body = nil
	ctx.responseEntity.contentType = DefaultMediaType
}

func (ctx *WebRequestContext) writeResponse() {
	ctx.fastHttpRequestContext.SetStatusCode(ctx.responseEntity.status)
	if ctx.responseEntity.contentType == MediaTypeApplicationJson {
		ctx.fastHttpRequestContext.SetContentType(MediaTypeApplicationJsonValue)

		if ctx.responseEntity.body == nil {
			return
		}
		result, err := json.Marshal(ctx.responseEntity.body)
		if err != nil {
			ctx.ThrowError(err)
		}
		ctx.fastHttpRequestContext.SetBody(result)
	} else if ctx.responseEntity.contentType == MediaTypeApplicationTextHtml {
		ctx.fastHttpRequestContext.SetContentType(MediaTypeApplicationTextHtmlValue)
		if ctx.responseEntity.body == nil {
			return
		}
		switch ctx.responseEntity.body.(type) {
		case string:
			value := []byte(ctx.responseEntity.body.(string))
			ctx.fastHttpRequestContext.SetBody(value)
		}
	} else {
		ctx.fastHttpRequestContext.SetContentType(MediaTypeApplicationXmlValue)

		if ctx.responseEntity.body == nil {
			return
		}

		result, err := xml.Marshal(ctx.responseEntity.body)
		if err != nil {
			ctx.ThrowError(err)
		}
		ctx.fastHttpRequestContext.SetBody(result)
	}
}

func (ctx *WebRequestContext) invoke(recoveryActive bool) {
	if recoveryActive {
		defer recoveryFunction(ctx)
		ctx.Next()
	} else {
		ctx.Next()
	}
}

func (ctx *WebRequestContext) Next() {
	if ctx.handlerIndex > ctx.handlerChain.handlerIndex {
		return
	}
next:
	if ctx.handlerIndex > ctx.handlerChain.handlerEndIndex {
		return
	}
	ctx.handlerChain.allHandlers[ctx.handlerIndex](ctx)
	ctx.handlerIndex++
	if ctx.handlerIndex == ctx.handlerChain.afterCompletionStartIndex {
		ctx.writeResponse()
		ctx.completedFlow = true
	}
	goto next
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

func (ctx *WebRequestContext) getPathByteArray() []byte {
	if ctx.uri == nil {
		ctx.uri = ctx.fastHttpRequestContext.URI()
		ctx.path = ctx.uri.Path()
	}
	return ctx.path
}

func (ctx *WebRequestContext) GetPath() string {
	if len(ctx.path) == 0 {
		return string(ctx.getPathByteArray())
	}
	return string(ctx.path)
}

func (ctx *WebRequestContext) GetPathVariable(name string) (string, bool) {
	for _, pathVariableName := range ctx.handlerChain.pathVariables {
		if pathVariableName == name {

		}
	}
	return "", false
}

func (ctx *WebRequestContext) GetRequestParameter(name string) (string, bool) {
	if ctx.args == nil {
		ctx.args = ctx.fastHttpRequestContext.QueryArgs()
	}
	result := ctx.args.Peek(name)
	if result == nil {
		return "", false
	}
	return string(result), true
}

func (ctx *WebRequestContext) GetHeaderValue(key string) (string, bool) {
	val := ctx.fastHttpRequestContext.Request.Header.Peek(key)
	if val == nil {
		return "", false
	}
	return string(val), true
}

func (ctx *WebRequestContext) GetRequest(request interface{}) {
	typ := reflect.TypeOf(request)
	if typ == nil {
		panic("Type cannot be determined as the given object is nil")
	}
	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
	}

	var cacheType *RequestObjectCache
	cacheRequestObjectMu.Lock()
	if cache, ok := cacheRequestObject[typ]; ok {
		cacheType = cache
		cacheRequestObjectMu.Unlock()
	} else {
		cacheRequestObjectMu.Unlock()
		return
	}

	body := ctx.fastHttpRequestContext.Request.Body()
	if cacheType.hasOnlyBody {
		contentType := core.BytesToStr(ctx.fastHttpRequestContext.Request.Header.Peek("Content-Type"))
		if contentType == MediaTypeApplicationJsonValue {
			err := json.Unmarshal(body, request)
			if err != nil {
				ctx.ThrowError(err)
			}
		} else {
			err := xml.Unmarshal(body, request)
			if err != nil {
				ctx.ThrowError(err)
			}
		}
		return
	}

	val := reflect.ValueOf(request)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	if cacheType.bodyFieldIndex != -1 {
		bodyValue := val.Field(cacheType.bodyFieldIndex)
		contentType := core.BytesToStr(ctx.fastHttpRequestContext.Request.Header.Peek("Content-Type"))
		if contentType == MediaTypeApplicationJsonValue {
			err := json.Unmarshal(body, bodyValue.Addr().Interface())
			if err != nil {
				ctx.ThrowError(err)
			}
		} else if contentType == MediaTypeApplicationXmlValue {
			err := xml.Unmarshal(body, bodyValue.Addr().Interface())
			if err != nil {
				ctx.ThrowError(err)
			}
		}
	}

}

func (ctx *WebRequestContext) SetStatus(status int) ResponseBodyBuilder {
	ctx.responseEntity.status = status
	return ctx
}

func (ctx *WebRequestContext) SetBody(body interface{}) ResponseBodyBuilder {
	if body == nil {
		return ctx
	}
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
	return ctx.err
}

func (ctx *WebRequestContext) SetError(err error) {
	ctx.err = err
}

func (ctx *WebRequestContext) ThrowError(err error) {
	panic(err)
}
