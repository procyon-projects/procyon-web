package web

type HandlerInvoker interface {
	InvokeHandler() error
}

type DefaultHandlerInvoker struct {
	invokerProcessor HandlerInvokerProcessor
}

func NewDefaultHandlerInvoker(invokerProcessor HandlerInvokerProcessor) DefaultHandlerInvoker {
	return DefaultHandlerInvoker{
		invokerProcessor,
	}
}

func (handlerInvoker DefaultHandlerInvoker) InvokeHandler() error {
	invokerProcessor := handlerInvoker.invokerProcessor
	if invokerProcessor == nil {
		return nil
	}
	_, err := invokerProcessor.PreProcess(nil)
	if err != nil {
		return err
	}
	_, err = invokerProcessor.Process(nil)
	if err != nil {
		return err
	}
	_, err = invokerProcessor.PostProcess(nil)
	if err != nil {
		return err
	}
	return nil
}

type HandlerInvokerProcessor interface {
	PreProcess(ctx ConfigurableWebApplicationContext) (ConfigurableWebApplicationContext, error)
	Process(ctx ConfigurableWebApplicationContext) (ConfigurableWebApplicationContext, error)
	PostProcess(ctx ConfigurableWebApplicationContext) (ConfigurableWebApplicationContext, error)
}
