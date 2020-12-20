package web

import core "github.com/procyon-projects/procyon-core"

type HandlerInterceptor func(requestContext *WebRequestContext)

type HandlerInterceptorBefore interface {
	HandleBefore(requestContext *WebRequestContext)
}

type HandlerInterceptorAfter interface {
	HandleAfter(requestContext *WebRequestContext)
}

type HandlerInterceptorAfterCompletion interface {
	AfterCompletion(requestContext *WebRequestContext)
}

type handlerInterceptorData struct {
	interceptorFunction HandlerInterceptor
	priority            core.PriorityValue
}

func newHandlerInterceptorData(interceptorFunction HandlerInterceptor, priority core.PriorityValue) *handlerInterceptorData {
	return &handlerInterceptorData{
		interceptorFunction: interceptorFunction,
		priority:            priority,
	}
}

type HandlerInterceptorRegistry interface {
	RegisterHandlerInterceptor(interceptorInstance interface{})
}

type SimpleHandlerInterceptorRegistry struct {
	beforeInterceptors          []*handlerInterceptorData
	afterInterceptors           []*handlerInterceptorData
	afterCompletionInterceptors []*handlerInterceptorData
}

func NewSimpleHandlerInterceptorRegistry() SimpleHandlerInterceptorRegistry {
	return SimpleHandlerInterceptorRegistry{
		beforeInterceptors:          make([]*handlerInterceptorData, 0),
		afterInterceptors:           make([]*handlerInterceptorData, 0),
		afterCompletionInterceptors: make([]*handlerInterceptorData, 0),
	}
}

func (registry SimpleHandlerInterceptorRegistry) RegisterHandlerInterceptor(interceptor interface{}) {
	priority := core.PriorityLowest
	if obj, ok := interceptor.(core.Priority); ok {
		priority = obj.GetPriority()
	}

	if interceptor, ok := interceptor.(HandlerInterceptorBefore); ok {
		registry.registerHandlerInterceptorBefore(priority, interceptor.HandleBefore)
	}

	if interceptor, ok := interceptor.(HandlerInterceptorAfter); ok {
		registry.registerHandlerInterceptorAfter(priority, interceptor.HandleAfter)
	}

	if interceptor, ok := interceptor.(HandlerInterceptorAfterCompletion); ok {
		registry.registerHandlerInterceptorAfterCompletion(priority, interceptor.AfterCompletion)
	}
}

func (registry SimpleHandlerInterceptorRegistry) registerHandlerInterceptorBefore(priority core.PriorityValue,
	interceptor HandlerInterceptor) {
	interceptorIndex := 0
	for index, registeredInterceptor := range registry.beforeInterceptors {
		if registeredInterceptor.priority > priority {
			interceptorIndex = index
		}
	}

	registry.beforeInterceptors = append(registry.beforeInterceptors, nil)
	copy(registry.beforeInterceptors[interceptorIndex+1:], registry.beforeInterceptors[interceptorIndex:])
	registry.beforeInterceptors[interceptorIndex] = newHandlerInterceptorData(interceptor, priority)
}

func (registry SimpleHandlerInterceptorRegistry) registerHandlerInterceptorAfter(priority core.PriorityValue,
	interceptor HandlerInterceptor) {
	interceptorIndex := 0
	for index, registeredInterceptor := range registry.afterInterceptors {
		if registeredInterceptor.priority > priority {
			interceptorIndex = index
		}
	}

	registry.afterInterceptors = append(registry.afterInterceptors, nil)
	copy(registry.afterInterceptors[interceptorIndex+1:], registry.afterInterceptors[interceptorIndex:])
	registry.afterInterceptors[interceptorIndex] = newHandlerInterceptorData(interceptor, priority)
}

func (registry SimpleHandlerInterceptorRegistry) registerHandlerInterceptorAfterCompletion(priority core.PriorityValue,
	interceptor HandlerInterceptor) {
	interceptorIndex := 0
	for index, registeredInterceptor := range registry.afterCompletionInterceptors {
		if registeredInterceptor.priority > priority {
			interceptorIndex = index
		}
	}

	registry.afterCompletionInterceptors = append(registry.afterCompletionInterceptors, nil)
	copy(registry.afterCompletionInterceptors[interceptorIndex+1:], registry.afterCompletionInterceptors[interceptorIndex:])
	registry.afterCompletionInterceptors[interceptorIndex] = newHandlerInterceptorData(interceptor, priority)
}
