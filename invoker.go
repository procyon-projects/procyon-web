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
	if handlerInvoker.invokerProcessor == nil {
		return nil
	}
	_, err := handlerInvoker.invokerProcessor.PreProcess(nil)
	if err != nil {
		return err
	}
	_, err = handlerInvoker.invokerProcessor.Process(nil)
	if err != nil {
		return err
	}
	_, err = handlerInvoker.invokerProcessor.PostProcess(nil)
	if err != nil {
		return err
	}
	return nil
}

type HandlerInvokerProcessor interface {
	PreProcess(ctx ConfigurableApplicationContext) (ConfigurableApplicationContext, error)
	Process(ctx ConfigurableApplicationContext) (ConfigurableApplicationContext, error)
	PostProcess(ctx ConfigurableApplicationContext) (ConfigurableApplicationContext, error)
}
