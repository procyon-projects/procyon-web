package web

import (
	"github.com/google/uuid"
	context "github.com/procyon-projects/procyon-context"
	peas "github.com/procyon-projects/procyon-peas"
)

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
	newContext := webTransactionContextPool.Get().(*BaseWebApplicationContext)
	newContext.BaseApplicationContext = ctx.(*context.BaseApplicationContext)
	newContext.SetParentPeaFactory(peaFactory)
	ctx.Copy(newContext, contextId)
	return newContext
}
