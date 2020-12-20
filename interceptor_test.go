package web

import (
	core "github.com/procyon-projects/procyon-core"
	"github.com/stretchr/testify/assert"
	"testing"
)

type testInterceptor1 struct {
}

func (interceptor testInterceptor1) HandleBefore(requestContext *WebRequestContext) {

}

func (interceptor testInterceptor1) HandleAfter(requestContext *WebRequestContext) {

}

func (interceptor testInterceptor1) AfterCompletion(requestContext *WebRequestContext) {

}

type testInterceptor2 struct {
}

func (interceptor testInterceptor2) HandleBefore(requestContext *WebRequestContext) {

}

func (interceptor testInterceptor2) GetPriority() core.PriorityValue {
	return core.PriorityHighest
}

func TestHandlerInterceptorRegistry_RegisterHandlerInterceptor(t *testing.T) {
	registry := NewSimpleHandlerInterceptorRegistry()
	registry.RegisterHandlerInterceptor(testInterceptor1{})
	assert.Len(t, registry.beforeInterceptors, 1)
	assert.Len(t, registry.afterInterceptors, 1)
	assert.Len(t, registry.afterCompletionInterceptors, 1)

	registry.RegisterHandlerInterceptor(testInterceptor2{})
	assert.Len(t, registry.beforeInterceptors, 2)
	assert.Len(t, registry.afterInterceptors, 1)
	assert.Len(t, registry.afterCompletionInterceptors, 1)
}
