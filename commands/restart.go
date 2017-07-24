package commands

import (
	"github.com/kpango/glg"
	"github.com/vetcher/easystarter/backend"
)

type RestartCommand struct {
}

func (c RestartCommand) Exec(args ...string) {
	if len(args) > 0 {
		svcName := args[0]
		if svcName == "-all" {
			backend.RestartServices(args...)
		} else {
			backend.RestartService(svcName)
		}
	} else {
		glg.Error("Specify service name.")
	}
}
