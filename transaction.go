package web

import (
	"errors"
	"github.com/google/uuid"
	context "github.com/procyon-projects/procyon-context"
	core "github.com/procyon-projects/procyon-core"
	peas "github.com/procyon-projects/procyon-peas"
	tx "github.com/procyon-projects/procyon-tx"
)

type TransactionContext struct {
	context.ConfigurableApplicationContext
	tx.TransactionalContext
}

func newTransactionContext() interface{} {
	return &TransactionContext{}
}

func prepareTransactionContext(contextId uuid.UUID,
	configurableContext context.ConfigurableApplicationContext,
	logger core.Logger) (*TransactionContext, error) {
	peaFactory, err := clonePeaFactoryForTransactionContext(configurableContext.GetPeaFactory())
	if err != nil {
		return nil, err
	}
	cloneContext := cloneApplicationContext(contextId, configurableContext, peaFactory, logger)
	var transactionalContext tx.TransactionalContext
	transactionalContext, err = nil, nil //tx.NewSimpleTransactionalContext(contextId, logger, nil, nil)
	if err != nil {
		return nil, err
	}
	txContext := transactionContextPool.Get().(*TransactionContext)
	if ctx, ok := cloneContext.(context.ConfigurableApplicationContext); ok {
		txContext.ConfigurableApplicationContext = ctx
		txContext.TransactionalContext = transactionalContext
	} else {
		transactionContextPool.Put(txContext)
		return nil, errors.New("context.ConfigurableApplicationContext methods must be implemented in your context struct")
	}
	return txContext, nil
}

func clonePeaFactoryForTransactionContext(parent peas.ConfigurablePeaFactory) (peas.ConfigurablePeaFactory, error) {
	peaFactory := parent.ClonePeaFactory().(peas.ConfigurablePeaFactory)
	peaFactory.SetParentPeaFactory(peaFactory)

	/* register scopes */
	err := peaFactory.RegisterScope(RequestScope, NewAppRequestScope())
	if err != nil {
		return nil, err
	}

	return peaFactory, nil
}

func cloneApplicationContext(contextId uuid.UUID,
	context context.ConfigurableApplicationContext,
	peaFactory peas.ConfigurablePeaFactory,
	logger core.Logger) interface{} {
	cloneContext := context.CloneContext(contextId, peaFactory)
	cloneContext.SetLogger(logger)
	return cloneContext
}
