package web

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

type TestController struct {
}

func (controller *TestController) HandlerFunc(testArg string) {
}

var (
	testController = &TestController{}
)

func TestWithGet(t *testing.T) {
	handlerMethod := WithGet("/test", testController.HandlerFunc)
	assert.Equal(t, "/test", handlerMethod.Path)
	assert.Equal(t, HttpMethodGet, handlerMethod.Method)
}

func TestWithPost(t *testing.T) {
	handlerMethod := WithPost("/test", testController.HandlerFunc)
	assert.Equal(t, "/test", handlerMethod.Path)
	assert.Equal(t, HttpMethodPost, handlerMethod.Method)
}

func TestWithPut(t *testing.T) {
	handlerMethod := WithPut("/test", testController.HandlerFunc)
	assert.Equal(t, "/test", handlerMethod.Path)
	assert.Equal(t, HttpMethodPut, handlerMethod.Method)
}

func TestWithDelete(t *testing.T) {
	handlerMethod := WithDelete("/test", testController.HandlerFunc)
	assert.Equal(t, "/test", handlerMethod.Path)
	assert.Equal(t, HttpMethodDelete, handlerMethod.Method)
}

func TestWithPatch(t *testing.T) {
	handlerMethod := WithPatch("/test", testController.HandlerFunc)
	assert.Equal(t, "/test", handlerMethod.Path)
	assert.Equal(t, HttpMethodPatch, handlerMethod.Method)
}

func TestHandlerMethodRegistry(t *testing.T) {
	registry := NewHandlerMethodRegistry()
	registry.Register(WithGet("/api/test", testController.HandlerFunc))
	assert.Equal(t, len(registry.registerMap), 1)
	registry.RegisterGroup("/api/test",
		WithGet("/", testController.HandlerFunc),
		WithPost("/", testController.HandlerFunc),
		WithDelete("/{id}", testController.HandlerFunc),
	)
	assert.Equal(t, len(registry.registerMap), 2)
	assert.Equal(t, len(registry.registerMap["/api/test"]), 3)
}
