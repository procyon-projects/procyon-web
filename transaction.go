package web

import (
	"github.com/google/uuid"
	context "github.com/procyon-projects/procyon-context"
	core "github.com/procyon-projects/procyon-core"
	peas "github.com/procyon-projects/procyon-peas"
	tx "github.com/procyon-projects/procyon-tx"
)

type TransactionContext struct {
	context.ConfigurableApplicationContext
	transactionalContext tx.TransactionalContext
}

func prepareTransactionContext(contextId uuid.UUID,
	context context.ConfigurableApplicationContext,
	logger core.Logger) (*TransactionContext, error) {
	peaFactory, err := clonePeaFactoryForTransactionContext(context.GetPeaFactory())
	if err != nil {
		return nil, err
	}
	_ = cloneApplicationContext(contextId, context, peaFactory, logger)
	var transactionalContext tx.TransactionalContext
	transactionalContext, err = tx.NewSimpleTransactionalContext(nil, nil)
	if err != nil {
		return nil, err
	}
	txContext := &TransactionContext{
		context,
		transactionalContext,
	}
	return txContext, nil
}

func clonePeaFactoryForTransactionContext(parent peas.ConfigurablePeaFactory) (peas.ConfigurablePeaFactory, error) {
	peaFactory := parent.ClonePeaFactory().(peas.ConfigurablePeaFactory)
	peaFactory.SetParentPeaFactory(peaFactory)

	/* register scopes */
	requestScope := NewAppRequestScope()
	err := peaFactory.RegisterScope(RequestScope, requestScope)
	if err != nil {
		return nil, err
	}
	/* register controller types to scope */
	err = peaFactory.RegisterTypeToScope(core.GetType((*Controller)(nil)), requestScope)
	if err != nil {
		return nil, err
	}
	err = peaFactory.RegisterTypeToScope(core.GetType((*context.Repository)(nil)), requestScope)
	if err != nil {
		return nil, err
	}
	err = peaFactory.RegisterTypeToScope(core.GetType((*context.Service)(nil)), requestScope)
	if err != nil {
		return nil, err
	}
	return peaFactory, nil
}

func cloneApplicationContext(contextId uuid.UUID,
	context context.ConfigurableApplicationContext,
	peaFactory peas.ConfigurablePeaFactory,
	logger core.Logger) context.ConfigurableContext {
	cloneContext := context.CloneContext(contextId, peaFactory)
	cloneContext.SetLogger(logger)
	return cloneContext
}
