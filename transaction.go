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
	logger               core.Logger
}

func prepareTransactionContext(context context.ConfigurableApplicationContext) *TransactionContext {
	_, err := preparePeaFactoryForTransactionContext(context.GetPeaFactory())
	if err != nil {

	}
	var transactionalContext tx.TransactionalContext
	transactionalContext, err = tx.NewSimpleTransactionalContext(nil, nil)
	if err != nil {

	}
	txContext := &TransactionContext{
		context,
		transactionalContext,
		nil,
	}
	// configure logger
	txContext.configureLogger()
	return txContext
}

func preparePeaFactoryForTransactionContext(parent peas.ConfigurablePeaFactory) (peas.ConfigurablePeaFactory, error) {
	peaFactory := parent.Clone().(peas.ConfigurablePeaFactory)
	peaFactory.SetParentPeaFactory(peaFactory)

	/* register scopes */
	err := peaFactory.RegisterScope(RequestScope, NewAppRequestScope())
	if err != nil {
		return nil, err
	}

	return peaFactory, nil
}

func (ctx TransactionContext) GetContextId() uuid.UUID {
	return ctx.transactionalContext.GetContextId()
}

func (ctx TransactionContext) GetLogger() core.Logger {
	return ctx.logger
}

func (ctx TransactionContext) configureLogger() {

}
