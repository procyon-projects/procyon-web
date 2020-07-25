package web

import core "github.com/procyon-projects/procyon-core"

func init() {
	/* Register Pool Types */
	registerPool(httpRequestType, newHttpRequest)
	registerPool(httpResponseType, newHttpResponse)
	/* Handler Info Registry */
	core.Register(NewSimpleHandlerInfoRegistry)
}
