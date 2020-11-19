package web

import (
	"github.com/procyon-projects/procyon-core"
)

type StandardWebEnvironment struct {
	core.StandardEnvironment
}

func NewStandardWebEnvironment() *StandardWebEnvironment {
	return &StandardWebEnvironment{
		core.NewStandardEnvironment(),
	}
}
