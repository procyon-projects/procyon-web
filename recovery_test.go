package web

import (
	"testing"
)

func Test_recovery(t *testing.T) {
	/*	webRequestContext := &WebRequestContext{}

		panicString(webRequestContext)
		assert.Equal(t, "error message", webRequestContext.err.Error())

		panicErr(webRequestContext)
		assert.Equal(t, "error message", webRequestContext.err.Error())

		panicUnknown(webRequestContext)
		assert.Equal(t, "unknown error", webRequestContext.err.Error())
	*/
}

func panicString(requestContext *WebRequestContext) {
	//defer recoveryFunction(requestContext)
	//panic("error message")
}

func panicErr(requestContext *WebRequestContext) {
	//defer recoveryFunction(requestContext)
	//panic(errors.New("error message"))
}

func panicUnknown(requestContext *WebRequestContext) {
	//defer recoveryFunction(requestContext)
	//panic(1)
}
