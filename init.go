package web

import core "github.com/procyon-projects/procyon-core"

func init() {
	/* Initialize Pools */
	initHttpRequestPool()
	initHttpResponsePool()
	initTransactionContextPool()
	/* Handler Info Registry */
	core.Register(NewSimpleHandlerInfoRegistry)
}
