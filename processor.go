package web

type RequestHandlerMappingProcessor struct {
	pathMatcher           PathMatcher
	requestHandlerMapping RequestHandlerMapping
}

func NewRequestHandlerMappingProcessor(pathMatcher PathMatcher, mapping RequestHandlerMapping) RequestHandlerMappingProcessor {
	return RequestHandlerMappingProcessor{
		pathMatcher,
		mapping,
	}
}

func (processor RequestHandlerMappingProcessor) AfterProperties() {

}

func (processor RequestHandlerMappingProcessor) BeforePeaInitialization(peaName string, pea interface{}) (interface{}, error) {
	if pea == nil {
		return nil, nil
	}
	if controller, ok := pea.(Controller); ok {
		handlerRegistry := newSimpleHandlerRegistry()
		controller.RegisterHandlers(handlerRegistry)
		processor.processHandler(peaName, handlerRegistry)
	}
	return pea, nil
}

func (processor RequestHandlerMappingProcessor) Initialize() error {
	return nil
}

func (processor RequestHandlerMappingProcessor) AfterPeaInitialization(peaName string, pea interface{}) (interface{}, error) {
	return pea, nil
}

func (processor RequestHandlerMappingProcessor) processHandler(handlerName string, handlerRegistry HandlerRegistry) {
	if handlerRegistry == nil {
		return
	}
	if simpleRegistry, ok := handlerRegistry.(SimpleHandlerRegistry); ok {
		registryMap := simpleRegistry.getRegistryMap()
		for prefix, handlers := range registryMap {
			for _, handler := range handlers {
				requestMappingInfo := processor.createRequestMapping(prefix, handler)
				processor.requestHandlerMapping.RegisterHandlerMethod(handlerName, requestMappingInfo, handler.HandlerFunc)
			}
		}
	}
}

func (processor RequestHandlerMappingProcessor) createRequestMapping(prefix string, handler RequestHandler) *RequestMapping {
	return NewRequestMapping(newMethodRequestMatcher(handler.Methods),
		newParametersRequestMatcher(),
		newPatternRequestMatcher(processor.pathMatcher, prefix, handler.Paths),
	)
}
