package commands

import (
	"github.com/kpango/glg"
	"github.com/vetcher/easystarter/backend"
)

type StopCommand struct {
}

func (c StopCommand) Exec(args ...string) error {
	if len(args) > 0 {
		svcName := args[0]
		if svcName == "-all" {
			backend.StopAllServices()
		} else {
			backend.StopService(svcName)
		}
	} else {
		glg.Error("Specify service name.")
	}
	return nil
}
