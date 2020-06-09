package web

import core "github.com/procyon-projects/procyon-core"

func init() {
	/* Handler Info Registry */
	core.Register(NewSimpleHandlerInfoRegistry)
}
