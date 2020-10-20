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
	return pea, nil
}

func (processor RequestHandlerMappingProcessor) Initialize() error {
	return nil
}

func (processor RequestHandlerMappingProcessor) AfterPeaInitialization(peaName string, pea interface{}) (interface{}, error) {
	return pea, nil
}
