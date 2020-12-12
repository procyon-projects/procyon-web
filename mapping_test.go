package web

import (
	"github.com/stretchr/testify/assert"
	"github.com/valyala/fasthttp"
	"testing"
)

func TestRequestMappingRegistry(t *testing.T) {
	handlerChain := NewHandlerChain(handlerFunction)
	registry := NewRequestMappingRegistry()
	registry.Register("/test", RequestMethodGet, handlerChain)

	context := newWebRequestContext().(*WebRequestContext)

	req := fasthttp.AcquireRequest()
	req.SetRequestURI("/test")
	req.Header.SetContentType("application/json")

	fastHttpRequestContext := &fasthttp.RequestCtx{}
	fastHttpRequestContext.Request = *req

	context.fastHttpRequestContext = fastHttpRequestContext

	registry.Find(context)
	assert.NotNil(t, context.handlerChain)
}

func TestRequestHandlerMapping(t *testing.T) {
	handlerMapping := NewRequestHandlerMapping(NewRequestMappingRegistry())
	handlerMapping.RegisterHandlerMethod("/test", RequestMethodGet, handlerFunction)

	context := newWebRequestContext().(*WebRequestContext)

	req := fasthttp.AcquireRequest()
	req.SetRequestURI("/test")
	req.Header.SetContentType("application/json")

	fastHttpRequestContext := &fasthttp.RequestCtx{}
	fastHttpRequestContext.Request = *req

	context.fastHttpRequestContext = fastHttpRequestContext

	handlerMapping.GetHandlerChain(context)
	assert.NotNil(t, context.handlerChain)
}
