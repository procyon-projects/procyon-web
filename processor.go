package web

type RequestHandlerMappingProcessor struct {
	requestHandlerMapping RequestHandlerMapping
}

func NewRequestHandlerMappingProcessor(mapping RequestHandlerMapping) RequestHandlerMappingProcessor {
	return RequestHandlerMappingProcessor{
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
		handlerRegistry := NewSimpleHandlerRegistry()
		controller.RegisterHandlers(handlerRegistry)
		processor.processHandler(handlerRegistry)
	}
	return pea, nil
}

func (processor RequestHandlerMappingProcessor) AfterPeaInitialization(peaName string, pea interface{}) (interface{}, error) {
	return pea, nil
}

func (processor RequestHandlerMappingProcessor) processHandler(handlerRegistry HandlerRegistry) {
	if handlerRegistry == nil {
		return
	}
	if simpleRegistry, ok := handlerRegistry.(SimpleHandlerRegistry); ok {
		registryMap := simpleRegistry.getRegistryMap()
		for prefix, handlers := range registryMap {
			for _, handler := range handlers {
				processor.requestHandlerMapping.RegisterHandlerMethod(prefix+handler.Path, handler.Method, handler.HandlerFunc)
			}
		}
	}
}
