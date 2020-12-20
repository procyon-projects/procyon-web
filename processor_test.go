package web

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

type testController struct {
}

func (controller testController) RegisterHandlers(registry HandlerRegistry) {
	registry.Register(Get(controller.handle, Path("/test")))
}

func (controller testController) handle(ctx *WebRequestContext) {

}

func TestRequestHandlerMappingProcessor(t *testing.T) {
	handlerMapping := NewRequestHandlerMapping(NewRequestMappingRegistry(), nil)
	processor := NewRequestHandlerMappingProcessor(handlerMapping)

	pea, err := processor.BeforePeaInitialization("", nil)
	assert.Nil(t, err)
	assert.Nil(t, pea)

	pea, err = processor.BeforePeaInitialization("testController", testController{})
	assert.Nil(t, err)
	assert.NotNil(t, pea)

	pea, err = processor.AfterPeaInitialization("testController", testController{})
	assert.Nil(t, err)
	assert.NotNil(t, pea)
}
