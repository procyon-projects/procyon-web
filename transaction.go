package web

import (
	"github.com/google/uuid"
	context "github.com/procyon-projects/procyon-context"
	peas "github.com/procyon-projects/procyon-peas"
)

func newApplicationContext() interface{} {
	return &context.BaseApplicationContext{}
}

func newWebTransactionContext() interface{} {
	return &BaseWebApplicationContext{}
}

func prepareWebTransactionContext(contextId uuid.UUID,
	configurableContext context.ConfigurableContext) (context.Context, error) {
	peaFactory, err := clonePeaFactory(configurableContext.GetPeaFactory())
	if err != nil {
		return nil, err
	}
	cloneContext := cloneWebTransactionContext(contextId, configurableContext, peaFactory)
	return cloneContext, nil
}

func clonePeaFactory(parent peas.ConfigurablePeaFactory) (peas.ConfigurablePeaFactory, error) {
	peaFactory := parent.ClonePeaFactory().(peas.ConfigurablePeaFactory)
	peaFactory.SetParentPeaFactory(peaFactory)

	/* register scopes */
	err := peaFactory.RegisterScope(RequestScope, NewAppRequestScope())
	if err != nil {
		return nil, err
	}

	return peaFactory, nil
}

func cloneWebTransactionContext(contextId uuid.UUID,
	ctx context.ConfigurableContext,
	peaFactory peas.ConfigurablePeaFactory) context.Context {

	newAppContext := applicationContextPool.Get().(*context.BaseApplicationContext)
	newAppContext.ConfigurableContextAdapter = ctx.(*BaseWebApplicationContext).ConfigurableContextAdapter
	newAppContext.ConfigurablePeaFactory = peaFactory
	ctx.Copy(newAppContext, contextId)

	newWebContext := webTransactionContextPool.Get().(*BaseWebApplicationContext)
	newWebContext.BaseApplicationContext = newAppContext

	return newWebContext
}
