package web

import (
	"github.com/procyon-projects/procyon-core"
)

type ConfigurableWebEnvironment interface {
	core.ConfigurableEnvironment
	Initialize()
}

type StandardWebEnvironment struct {
	core.StandardEnvironment
}

func NewStandardWebEnvironment() *StandardWebEnvironment {
	return &StandardWebEnvironment{
		core.NewStandardEnvironment(),
	}
}

func (env *StandardWebEnvironment) Initialize() {

}
