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
		processor.processHandler(handlerRegistry)
	}
	return pea, nil
}

func (processor RequestHandlerMappingProcessor) Initialize() error {
	return nil
}

func (processor RequestHandlerMappingProcessor) AfterPeaInitialization(peaName string, pea interface{}) (interface{}, error) {
	return pea, nil
}

func (processor RequestHandlerMappingProcessor) processHandler(handlerRegistry HandlerRegistry) {
	if handlerRegistry != nil {
		if simpleRegistry, ok := handlerRegistry.(*SimpleHandlerRegistry); ok {
			simpleRegistry.getRegistryMap()
		}
	}
}
