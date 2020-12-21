package web

import (
	"fmt"
	json "github.com/json-iterator/go"
	context "github.com/procyon-projects/procyon-context"
	"net/http"
)

var (
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
		handler.logger.Error(requestContext, err.Error())
		requestContext.SetStatus(httpError.Code)
	} else {
		handler.logger.Error(requestContext, err.Error())
		requestContext.SetStatus(http.StatusInternalServerError)
	}

	response, _ := json.Marshal(requestContext.responseEntity.body)
	requestContext.SetBody(response)
	requestContext.SetContentType(MediaTypeApplicationJson)
}
