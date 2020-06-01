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
	handler := NewHandlerInfo(testController.HandlerFunc, WithPath("/test"))
	assert.Equal(t, "/test", handler.Paths[0])
	assert.Equal(t, HttpMethodGet, handler.Methods[1])
}

func TestHandlerMethodRegistry(t *testing.T) {
	registry := NewSimpleHandlerInfoRegistry()
	registry.Register(NewHandlerInfo(testController.HandlerFunc, WithPath("/api/test")))
	assert.Equal(t, len(registry.infoRegistryMap), 1)
	registry.RegisterGroup("/api/test",
		NewHandlerInfo(testController.HandlerFunc, WithPath("/")),
		NewHandlerInfo(testController.HandlerFunc, WithPath("/"), WithMethod(HttpMethodPost)),
		NewHandlerInfo(testController.HandlerFunc, WithPath("/{id}"), WithMethod(HttpMethodDelete)),
	)
	assert.Equal(t, len(registry.infoRegistryMap), 2)
	assert.Equal(t, len(registry.infoRegistryMap["/api/test"]), 3)
}

func Test(t *testing.T) {
	x := NewProcyonServerApplicationContext()
	x.Configure()
	x.GetWebServer()
}
