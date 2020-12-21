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
		switch val := r.(type) {
		case string:
			ctx.err = errors.New(val)
		case error:
			ctx.err = val
		default:
			ctx.err = errors.New("unknown error : " + string(debug.Stack()))
		}
		recoveryManager.HandleError(ctx.err, ctx)
	}
}

func (recoveryManager *recoveryManager) HandleError(err error, ctx *WebRequestContext) {
	defer func() {
		if r := recover(); r != nil {
			switch err := r.(type) {
			case string:
				recoveryManager.logger.Error(ctx, err)
			case error:
				recoveryManager.logger.Error(ctx, err.Error())
			default:
				recoveryManager.logger.Error(ctx, "unknown error : "+string(debug.Stack()))
			}
		}
	}()

	if recoveryManager.customErrorHandler == nil {
		recoveryManager.customErrorHandler.HandleError(ctx.err, ctx)
	} else {
		recoveryManager.defaultErrorHandler.HandleError(ctx.err, ctx)
	}

	if ctx.handlerChain != nil {
		ctx.handlerIndex = ctx.handlerChain.afterCompletionStartIndex
		ctx.invokeHandlers(true)
	}
}
