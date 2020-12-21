package web

import (
	"errors"
	"github.com/procyon-projects/goo"
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
	ctx.responseEntity.body = "test-body"
	ctx.responseEntity.contentType = MediaTypeApplicationJson

	ctx.reset()

	assert.Equal(t, 0, ctx.handlerIndex)
	assert.Equal(t, 0, ctx.pathVariableCount)
	assert.Nil(t, ctx.valueMap)
	assert.Equal(t, http.StatusOK, ctx.responseEntity.status)
	assert.Nil(t, ctx.responseEntity.body)
	assert.Equal(t, DefaultMediaType, ctx.responseEntity.contentType)
}

func TestWebRequestContext_ValueMap(t *testing.T) {
	ctx := newWebRequestContext().(*WebRequestContext)
	ctx.Put("test-key", "test-value")
	assert.Equal(t, "test-value", ctx.Get("test-key"))
}

func TestWebRequestContext_Status(t *testing.T) {
	ctx := newWebRequestContext().(*WebRequestContext)
	ctx.SetStatus(http.StatusNotFound)
	assert.Equal(t, http.StatusNotFound, ctx.GetStatus())
}

func TestWebRequestContext_Body(t *testing.T) {
	ctx := newWebRequestContext().(*WebRequestContext)
	ctx.SetBody("test-body")
	assert.Equal(t, "test-body", ctx.GetBody())
}

func TestWebRequestContext_ContextType(t *testing.T) {
	ctx := newWebRequestContext().(*WebRequestContext)
	ctx.SetContentType(MediaTypeApplicationJson)
	assert.Equal(t, MediaTypeApplicationJson, ctx.GetContentType())
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
	ctx.SetContentType(MediaTypeApplicationTextHtml)
	ctx.SetBody("test")
	ctx.writeResponse()
	assert.Equal(t, "test", string(ctx.fastHttpRequestContext.Response.Body()))
	assert.Equal(t, MediaTypeApplicationTextHtmlValue, string(ctx.fastHttpRequestContext.Response.Header.ContentType()))
}

func TestWebRequestContext_writeResponseAsJson(t *testing.T) {
	ctx := newWebRequestContext().(*WebRequestContext)
	ctx.fastHttpRequestContext = &fasthttp.RequestCtx{}
	ctx.SetContentType(MediaTypeApplicationJson)
	ctx.SetBody(testResponse{"test", 25})
	ctx.writeResponse()
	assert.Equal(t, "{\"Name\":\"test\",\"Age\":25}", string(ctx.fastHttpRequestContext.Response.Body()))
	assert.Equal(t, MediaTypeApplicationJsonValue, string(ctx.fastHttpRequestContext.Response.Header.ContentType()))
}

func TestWebRequestContext_writeResponseAsXml(t *testing.T) {
	ctx := newWebRequestContext().(*WebRequestContext)
	ctx.fastHttpRequestContext = &fasthttp.RequestCtx{}
	ctx.SetContentType(MediaTypeApplicationXml)
	ctx.SetBody(testResponse{"test", 25})
	ctx.writeResponse()
	assert.Equal(t, "<testResponse><Name>test</Name><Age>25</Age></testResponse>", string(ctx.fastHttpRequestContext.Response.Body()))
	assert.Equal(t, MediaTypeApplicationXmlValue, string(ctx.fastHttpRequestContext.Response.Header.ContentType()))
}

func TestWebRequestContext_GetRequestWithNil(t *testing.T) {
	ctx := newWebRequestContext().(*WebRequestContext)
	assert.Panics(t, func() {
		ctx.GetRequest(nil)
	})
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
	scanRequestObject(goo.GetType(requestObj))
	ctx.GetRequest(requestObj)

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
	scanRequestObject(goo.GetType(requestObj))
	ctx.GetRequest(requestObj)

	assert.Equal(t, requestObj.Body.Name, "test")
	assert.Equal(t, requestObj.Body.Age, 25)
}

func TestWebRequestContext_GetRequestForJson_WithOnlyBody(t *testing.T) {
	ctx := newWebRequestContext().(*WebRequestContext)
	ctx.fastHttpRequestContext = &fasthttp.RequestCtx{}
	req := fasthttp.AcquireRequest()
	req.SetBody([]byte("{\"Name\":\"test\",\"Age\":25}"))
	req.Header.SetContentType(MediaTypeApplicationJsonValue)
	ctx.fastHttpRequestContext.Request = *req

	requestObj := &testRequestObjectWithOnlyBody{}
	scanRequestObject(goo.GetType(requestObj))
	ctx.GetRequest(requestObj)

	assert.Equal(t, requestObj.Name, "test")
	assert.Equal(t, requestObj.Age, 25)
}

func TestWebRequestContext_GetRequestForXml_WithOnlyBody(t *testing.T) {
	ctx := newWebRequestContext().(*WebRequestContext)
	ctx.fastHttpRequestContext = &fasthttp.RequestCtx{}
	req := fasthttp.AcquireRequest()
	req.SetBody([]byte("<testRequestObjectWithOnlyBody><Name>test</Name><Age>25</Age></testRequestObjectWithOnlyBody>"))
	req.Header.SetContentType(MediaTypeApplicationXmlValue)
	ctx.fastHttpRequestContext.Request = *req

	requestObj := &testRequestObjectWithOnlyBody{}
	scanRequestObject(goo.GetType(requestObj))
	ctx.GetRequest(requestObj)

	assert.Equal(t, requestObj.Name, "test")
	assert.Equal(t, requestObj.Age, 25)
}
