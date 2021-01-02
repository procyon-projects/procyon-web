package web

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/valyala/fasthttp"
	"net/http"
	"testing"
)

func TestWebRequestContext_prepare(t *testing.T) {
	ctx := newWebRequestContext().(*WebRequestContext)

	ctx.prepare(false)
	assert.Equal(t, 0, len(ctx.contextIdStr))

	ctx.prepare(true)
	assert.NotNil(t, 0, len(ctx.contextIdStr))
}

func TestWebRequestContext_reset(t *testing.T) {
	ctx := newWebRequestContext().(*WebRequestContext)

	ctx.handlerIndex = 1
	ctx.pathVariableCount = 1
	ctx.valueMap = make(map[string]interface{})
	ctx.responseEntity.status = http.StatusCreated
	ctx.responseEntity.model = "test-body"
	ctx.responseEntity.contentType = MediaTypeApplicationJson

	ctx.reset()

	assert.Equal(t, 0, ctx.handlerIndex)
	assert.Equal(t, 0, ctx.pathVariableCount)
	assert.Nil(t, ctx.valueMap)
	assert.Equal(t, http.StatusOK, ctx.responseEntity.status)
	assert.Nil(t, ctx.responseEntity.model)
	assert.Equal(t, DefaultMediaType, ctx.responseEntity.contentType)
}

func TestWebRequestContext_ValueMap(t *testing.T) {
	ctx := newWebRequestContext().(*WebRequestContext)
	ctx.Put("test-key", "test-value")
	assert.Equal(t, "test-value", ctx.Get("test-key"))
}

func TestWebRequestContext_Status(t *testing.T) {
	ctx := newWebRequestContext().(*WebRequestContext)
	ctx.SetResponseStatus(http.StatusNotFound)
	assert.Equal(t, http.StatusNotFound, ctx.GetResponseStatus())
}

func TestWebRequestContext_Body(t *testing.T) {
	ctx := newWebRequestContext().(*WebRequestContext)
	ctx.SetModel("test-body")
	assert.Equal(t, "test-body", ctx.GetModel())
}

func TestWebRequestContext_ContextType(t *testing.T) {
	ctx := newWebRequestContext().(*WebRequestContext)
	ctx.SetResponseContentType(MediaTypeApplicationJson)
	assert.Equal(t, MediaTypeApplicationJson, ctx.GetResponseContentType())
}

func TestWebRequestContext_Ok(t *testing.T) {
	ctx := newWebRequestContext().(*WebRequestContext)
	ctx.Ok()
	assert.Equal(t, http.StatusOK, ctx.responseEntity.status)
}

func TestWebRequestContext_NotFound(t *testing.T) {
	ctx := newWebRequestContext().(*WebRequestContext)
	ctx.NotFound()
	assert.Equal(t, http.StatusNotFound, ctx.responseEntity.status)
}

func TestWebRequestContext_NoContent(t *testing.T) {
	ctx := newWebRequestContext().(*WebRequestContext)
	ctx.NoContent()
	assert.Equal(t, http.StatusNoContent, ctx.responseEntity.status)
}

func TestWebRequestContext_BadRequest(t *testing.T) {
	ctx := newWebRequestContext().(*WebRequestContext)
	ctx.BadRequest()
	assert.Equal(t, http.StatusBadRequest, ctx.responseEntity.status)
}

func TestWebRequestContext_Accepted(t *testing.T) {
	ctx := newWebRequestContext().(*WebRequestContext)
	ctx.Accepted()
	assert.Equal(t, http.StatusAccepted, ctx.responseEntity.status)
}

func TestWebRequestContext_Created(t *testing.T) {
	ctx := newWebRequestContext().(*WebRequestContext)
	ctx.Created("")
	assert.Equal(t, http.StatusCreated, ctx.responseEntity.status)
}

func TestWebRequestContext_Error(t *testing.T) {
	//	ctx := newWebRequestContext().(*WebRequestContext)
	//	ctx.SetError(errors.New("test-error"))
	//	assert.Equal(t, "test-error", ctx.GetError().Error())
}

func TestWebRequestContext_ThrowError(t *testing.T) {
	ctx := newWebRequestContext().(*WebRequestContext)
	assert.Panics(t, func() {
		ctx.ThrowError(errors.New("test-error"))
	})
}

type testResponse struct {
	Name string
	Age  int
}

func TestWebRequestContext_writeResponseAsTextHtml(t *testing.T) {
	ctx := newWebRequestContext().(*WebRequestContext)
	ctx.fastHttpRequestContext = &fasthttp.RequestCtx{}
	ctx.SetResponseContentType(MediaTypeApplicationTextHtml)
	ctx.SetModel("test")
	ctx.writeResponse()
	assert.Equal(t, "test", string(ctx.fastHttpRequestContext.Response.Body()))
	assert.Equal(t, MediaTypeApplicationTextHtmlValue, string(ctx.fastHttpRequestContext.Response.Header.ContentType()))
}

func TestWebRequestContext_writeResponseAsJson(t *testing.T) {
	ctx := newWebRequestContext().(*WebRequestContext)
	ctx.fastHttpRequestContext = &fasthttp.RequestCtx{}
	ctx.SetResponseContentType(MediaTypeApplicationJson)
	ctx.SetModel(testResponse{"test", 25})
	ctx.writeResponse()
	assert.Equal(t, "{\"Name\":\"test\",\"Age\":25}", string(ctx.fastHttpRequestContext.Response.Body()))
	assert.Equal(t, MediaTypeApplicationJsonValue, string(ctx.fastHttpRequestContext.Response.Header.ContentType()))
}

func TestWebRequestContext_writeResponseAsXml(t *testing.T) {
	ctx := newWebRequestContext().(*WebRequestContext)
	ctx.fastHttpRequestContext = &fasthttp.RequestCtx{}
	ctx.SetResponseContentType(MediaTypeApplicationXml)
	ctx.SetModel(testResponse{"test", 25})
	ctx.writeResponse()
	assert.Equal(t, "<testResponse><Name>test</Name><Age>25</Age></testResponse>", string(ctx.fastHttpRequestContext.Response.Body()))
	assert.Equal(t, MediaTypeApplicationXmlValue, string(ctx.fastHttpRequestContext.Response.Header.ContentType()))
}

func TestWebRequestContext_GetRequestWithNil(t *testing.T) {
	ctx := newWebRequestContext().(*WebRequestContext)
	err := ctx.BindRequest(nil)
	assert.NotNil(t, err)
}

type testRequestObject struct {
	Body struct {
		Name string `json:"Name" yaml:"Name"`
		Age  int    `json:"Age" yaml:"Age"`
	} `request:"body"`
}

type testRequestObjectWithOnlyBody struct {
	Name string `json:"Name" yaml:"Name"`
	Age  int    `json:"Age" yaml:"Age"`
}

func TestWebRequestContext_GetRequestForJson(t *testing.T) {
	ctx := newWebRequestContext().(*WebRequestContext)
	ctx.fastHttpRequestContext = &fasthttp.RequestCtx{}
	req := fasthttp.AcquireRequest()
	req.SetBody([]byte("{\"Name\":\"test\",\"Age\":25}"))
	req.Header.SetContentType(MediaTypeApplicationJsonValue)
	ctx.fastHttpRequestContext.Request = *req

	requestObj := &testRequestObject{}
	ctx.handlerChain = NewHandlerChain(nil, nil, ScanRequestObjectMetadata(requestObj))
	ctx.BindRequest(requestObj)

	assert.Equal(t, requestObj.Body.Name, "test")
	assert.Equal(t, requestObj.Body.Age, 25)
}

func TestWebRequestContext_GetRequestForXml(t *testing.T) {
	ctx := newWebRequestContext().(*WebRequestContext)
	ctx.fastHttpRequestContext = &fasthttp.RequestCtx{}
	req := fasthttp.AcquireRequest()
	req.SetBody([]byte("<testRequestObject><Name>test</Name><Age>25</Age></testRequestObject>"))
	req.Header.SetContentType(MediaTypeApplicationXmlValue)
	ctx.fastHttpRequestContext.Request = *req

	requestObj := &testRequestObject{}
	ctx.handlerChain = NewHandlerChain(nil, nil, ScanRequestObjectMetadata(requestObj))
	ctx.BindRequest(requestObj)

	assert.Equal(t, requestObj.Body.Name, "test")
	assert.Equal(t, requestObj.Body.Age, 25)
}

func TestWebRequestContext_BindRequestForJson_WithOnlyBody(t *testing.T) {
	ctx := newWebRequestContext().(*WebRequestContext)
	ctx.fastHttpRequestContext = &fasthttp.RequestCtx{}
	req := fasthttp.AcquireRequest()
	req.SetBody([]byte("{\"Name\":\"test\",\"Age\":25}"))
	req.Header.SetContentType(MediaTypeApplicationJsonValue)
	ctx.fastHttpRequestContext.Request = *req

	requestObj := &testRequestObjectWithOnlyBody{}
	ctx.handlerChain = NewHandlerChain(nil, nil, ScanRequestObjectMetadata(requestObj))
	ctx.BindRequest(requestObj)

	assert.Equal(t, requestObj.Name, "test")
	assert.Equal(t, requestObj.Age, 25)
}

func TestWebRequestContext_BindRequestForXml_WithOnlyBody(t *testing.T) {
	ctx := newWebRequestContext().(*WebRequestContext)
	ctx.fastHttpRequestContext = &fasthttp.RequestCtx{}
	req := fasthttp.AcquireRequest()
	req.SetBody([]byte("<testRequestObjectWithOnlyBody><Name>test</Name><Age>25</Age></testRequestObjectWithOnlyBody>"))
	req.Header.SetContentType(MediaTypeApplicationXmlValue)
	ctx.fastHttpRequestContext.Request = *req

	requestObj := &testRequestObjectWithOnlyBody{}
	ctx.handlerChain = NewHandlerChain(nil, nil, ScanRequestObjectMetadata(requestObj))
	ctx.BindRequest(requestObj)

	assert.Equal(t, requestObj.Name, "test")
	assert.Equal(t, requestObj.Age, 25)
}
