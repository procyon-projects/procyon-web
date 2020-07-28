package web

import core "github.com/procyon-projects/procyon-core"

func init() {
	/* Initialize Pools */
	initHttpRequestPool()
	initHttpResponsePool()
	/* Handler Info Registry */
	core.Register(NewSimpleHandlerInfoRegistry)
}
