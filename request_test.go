package web

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

type requestObject struct {
}

func handlerFunction(ctx *WebRequestContext) {

}

func TestRequestHandler_Get(t *testing.T) {
	requestObject := requestObject{}
	requestHandler := Get(handlerFunction, Path("/test"), RequestObject(requestObject))
	assert.Equal(t, RequestMethodGet, requestHandler.Method)
	assert.Equal(t, "/test", requestHandler.Path)
	assert.Equal(t, requestObject, requestHandler.RequestObject)
}

func TestRequestHandler_Post(t *testing.T) {
	requestObject := requestObject{}
	requestHandler := Post(handlerFunction, Path("/test"), RequestObject(requestObject))
	assert.Equal(t, RequestMethodPost, requestHandler.Method)
	assert.Equal(t, "/test", requestHandler.Path)
	assert.Equal(t, requestObject, requestHandler.RequestObject)
}

func TestRequestHandler_Put(t *testing.T) {
	requestObject := requestObject{}
	requestHandler := Put(handlerFunction, Path("/test"), RequestObject(requestObject))
	assert.Equal(t, RequestMethodPut, requestHandler.Method)
	assert.Equal(t, "/test", requestHandler.Path)
	assert.Equal(t, requestObject, requestHandler.RequestObject)
}

func TestRequestHandler_Delete(t *testing.T) {
	requestObject := requestObject{}
	requestHandler := Delete(handlerFunction, Path("/test"), RequestObject(requestObject))
	assert.Equal(t, RequestMethodDelete, requestHandler.Method)
	assert.Equal(t, "/test", requestHandler.Path)
	assert.Equal(t, requestObject, requestHandler.RequestObject)
}

func TestRequestHandler_Patch(t *testing.T) {
	requestObject := requestObject{}
	requestHandler := Patch(handlerFunction, Path("/test"), RequestObject(requestObject))
	assert.Equal(t, RequestMethodPatch, requestHandler.Method)
	assert.Equal(t, "/test", requestHandler.Path)
	assert.Equal(t, requestObject, requestHandler.RequestObject)
}

func TestRequestHandler_Options(t *testing.T) {
	requestObject := requestObject{}
	requestHandler := Options(handlerFunction, Path("/test"), RequestObject(requestObject))
	assert.Equal(t, RequestMethodOptions, requestHandler.Method)
	assert.Equal(t, "/test", requestHandler.Path)
	assert.Equal(t, requestObject, requestHandler.RequestObject)
}

func TestRequestHandler_Head(t *testing.T) {
	requestObject := requestObject{}
	requestHandler := Head(handlerFunction, Path("/test"), RequestObject(requestObject))
	assert.Equal(t, RequestMethodHead, requestHandler.Method)
	assert.Equal(t, "/test", requestHandler.Path)
	assert.Equal(t, requestObject, requestHandler.RequestObject)
}

func TestSimpleHandlerRegistry_Register(t *testing.T) {
	requestObject := requestObject{}
	registry := NewSimpleHandlerRegistry()
	registry.Register(Get(handlerFunction, Path("/test"), RequestObject(requestObject)))
	assert.Equal(t, 1, len(registry.getRegistryMap()))
	assert.Equal(t, RequestMethodGet, registry.registryMap[""][0].Method)
	assert.Equal(t, "/test", registry.registryMap[""][0].Path)
	assert.Equal(t, requestObject, registry.registryMap[""][0].RequestObject)
}

func TestSimpleHandlerRegistry_RegisterGroup(t *testing.T) {
	requestObject := requestObject{}
	registry := NewSimpleHandlerRegistry()
	registry.RegisterGroup("/api",
		Get(handlerFunction, Path("/test"), RequestObject(requestObject)),
		Post(handlerFunction, Path("/test"), RequestObject(requestObject)),
	)
	assert.Equal(t, 1, len(registry.getRegistryMap()))
	assert.Equal(t, RequestMethodGet, registry.registryMap["/api"][0].Method)
	assert.Equal(t, "/test", registry.registryMap["/api"][0].Path)
	assert.Equal(t, requestObject, registry.registryMap["/api"][0].RequestObject)

	assert.Equal(t, RequestMethodPost, registry.registryMap["/api"][1].Method)
	assert.Equal(t, "/test", registry.registryMap["/api"][1].Path)
	assert.Equal(t, requestObject, registry.registryMap["/api"][1].RequestObject)

	registry = NewSimpleHandlerRegistry()
	assert.Panics(t, func() {
		registry.Register(Get(handlerFunction))
	})
}
