package web

type RequestHandlerMappingProcessor struct {
	requestHandlerMapping RequestHandlerMapping
}

func NewRequestHandlerMappingProcessor() RequestHandlerMappingProcessor {
	return RequestHandlerMappingProcessor{
		NewRequestHandlerMapping(),
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
				requestMappingInfo := processor.createRequestMappingInfo(prefix, handler)
				processor.requestHandlerMapping.RegisterHandlerMethod(handlerName, requestMappingInfo, handler.HandlerFunc)
			}
		}
	}
}

func (processor RequestHandlerMappingProcessor) createRequestMappingInfo(prefix string, handler RequestHandler) RequestMappingInfo {
	return newRequestMappingInfo("",
		newMethodRequestMatcher(handler.Methods),
		newParametersRequestMatcher(),
		newPatternRequestMatcher(prefix, handler.Paths),
	)
}
