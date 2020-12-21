package web

import (
	"errors"
	context "github.com/procyon-projects/procyon-context"
	"runtime/debug"
)

type recoveryManager struct {
	defaultErrorHandler ErrorHandler
	customErrorHandler  ErrorHandler
	logger              context.Logger
}

func newRecoveryManager(logger context.Logger) *recoveryManager {
	return &recoveryManager{
		defaultErrorHandler: NewDefaultErrorHandler(logger),
		logger:              logger,
	}
}

func (recoveryManager *recoveryManager) Recover(ctx *WebRequestContext) {
	if r := recover(); r != nil {
		ctx.crashed = true
		switch err := r.(type) {
		case string:
			ctx.err = errors.New(err)
		case error:
			ctx.err = err
		default:
			ctx.err = errors.New("unknown error : \n" + string(debug.Stack()))
		}
		recoveryManager.HandleError(ctx.err, ctx)
	}
}

func (recoveryManager *recoveryManager) HandleError(err error, ctx *WebRequestContext) {
	defer func() {
		if r := recover(); r != nil {

			var errText string
			switch err := r.(type) {
			case string:
				errText = err
			case error:
				errText = err.Error()
			default:
				errText = "unknown error : "
			}

			recoveryManager.logger.Error(ctx, errText+"\n"+string(debug.Stack()))
			if recoveryManager.customErrorHandler != nil {
				recoveryManager.defaultErrorHandler.HandleError(ctx.err, ctx)
				ctx.writeResponse()
			}
		}
	}()

	if recoveryManager.customErrorHandler != nil {
		recoveryManager.customErrorHandler.HandleError(err, ctx)
	} else {
		recoveryManager.defaultErrorHandler.HandleError(err, ctx)
	}
	ctx.writeResponse()

	if ctx.handlerChain != nil && ctx.handlerIndex < ctx.handlerChain.handlerIndex {
		ctx.handlerIndex = ctx.handlerChain.afterCompletionStartIndex
		ctx.invokeHandlers(true, nil)
	}
}
