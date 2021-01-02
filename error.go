package web

import (
	"errors"
	"fmt"
	context "github.com/procyon-projects/procyon-context"
	"net/http"
	"runtime/debug"
)

var (
	HttpErrorNoContent             = NewHTTPError(http.StatusNoContent)
	HttpErrorBadRequest            = NewHTTPError(http.StatusBadRequest)
	HttpErrorUnauthorized          = NewHTTPError(http.StatusUnauthorized)
	HttpErrorForbidden             = NewHTTPError(http.StatusForbidden)
	HttpErrorNotFound              = NewHTTPError(http.StatusNotFound)
	HttpErrorMethodNotAllowed      = NewHTTPError(http.StatusMethodNotAllowed)
	HttpErrorRequestTimeout        = NewHTTPError(http.StatusRequestTimeout)
	HttpErrorRequestEntityTooLarge = NewHTTPError(http.StatusRequestEntityTooLarge)
	HttpErrorUnsupportedMediaType  = NewHTTPError(http.StatusUnsupportedMediaType)
	HttpErrorTooManyRequests       = NewHTTPError(http.StatusTooManyRequests)

	HttpErrorInternalServerError = NewHTTPError(http.StatusInternalServerError)
	HttpErrorBadGateway          = NewHTTPError(http.StatusBadGateway)
	HttpErrorServiceUnavailable  = NewHTTPError(http.StatusServiceUnavailable)
)

func (err *HTTPError) Error() string {
	return fmt.Sprintf("code=%d, message=%v", err.Code, err.Message)
}

type HTTPError struct {
	Code    int
	Message interface{}
}

func NewHTTPError(code int, message ...interface{}) *HTTPError {
	httpError := &HTTPError{
		Code:    code,
		Message: http.StatusText(code),
	}

	if len(message) > 0 {
		httpError.Message = message[0]
	}

	return httpError
}

type ErrorHandler interface {
	HandleError(err error, requestContext *WebRequestContext)
}

type DefaultErrorHandler struct {
	logger context.Logger
}

func NewDefaultErrorHandler(logger context.Logger) DefaultErrorHandler {
	return DefaultErrorHandler{
		logger,
	}
}

func (handler DefaultErrorHandler) HandleError(err error, requestContext *WebRequestContext) {
	if httpError, ok := err.(*HTTPError); ok {
		requestContext.SetResponseStatus(httpError.Code)
		requestContext.SetModel(httpError)
	} else {
		handler.logger.Error(requestContext, err.Error()+"\n"+string(debug.Stack()))
		requestContext.SetResponseStatus(HttpErrorInternalServerError.Code)
		requestContext.SetModel(HttpErrorInternalServerError)
	}

	requestContext.SetResponseContentType(MediaTypeApplicationJson)
}

type errorHandlerManager struct {
	defaultErrorHandler ErrorHandler
	customErrorHandler  ErrorHandler
	logger              context.Logger
}

func newErrorHandlerManager(logger context.Logger) *errorHandlerManager {
	return &errorHandlerManager{
		defaultErrorHandler: NewDefaultErrorHandler(logger),
		logger:              logger,
	}
}

func (errorHandlerManager *errorHandlerManager) Recover(ctx *WebRequestContext) {
	if r := recover(); r != nil {
		ctx.crashed = true
		switch err := r.(type) {
		case *HTTPError:
			ctx.httpError = err
			errorHandlerManager.HandleError(ctx.httpError, ctx)
			return
		case string:
			ctx.internalError = errors.New(err)
		case error:
			ctx.internalError = err
		default:
			ctx.internalError = errors.New("unknown error : \n" + string(debug.Stack()))
		}
		errorHandlerManager.HandleError(ctx.internalError, ctx)
	}
}

func (errorHandlerManager *errorHandlerManager) JustHandleError(err error, ctx *WebRequestContext) {
	if errorHandlerManager.customErrorHandler != nil {
		errorHandlerManager.customErrorHandler.HandleError(err, ctx)
	} else {
		errorHandlerManager.defaultErrorHandler.HandleError(err, ctx)
	}
}

func (errorHandlerManager *errorHandlerManager) HandleError(err error, ctx *WebRequestContext) {
	defer errorHandlerManager.wtf(err, ctx)

	if errorHandlerManager.customErrorHandler != nil {
		errorHandlerManager.customErrorHandler.HandleError(err, ctx)
	} else {
		errorHandlerManager.defaultErrorHandler.HandleError(err, ctx)
	}
	ctx.writeResponse()

	if ctx.handlerChain != nil && ctx.handlerIndex < ctx.handlerChain.handlerIndex {
		ctx.handlerIndex = ctx.handlerChain.afterCompletionStartIndex
		ctx.invokeHandlers()
	}
}

func (errorHandlerManager *errorHandlerManager) wtf(err error, ctx *WebRequestContext) {
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

		errorHandlerManager.logger.Error(ctx, errText+"\n"+string(debug.Stack()))
		if errorHandlerManager.customErrorHandler != nil {
			errorHandlerManager.defaultErrorHandler.HandleError(err, ctx)
			ctx.writeResponse()
		}
	}
}
