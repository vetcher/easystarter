package commands

import (
	"github.com/kpango/glg"
	"github.com/vetcher/easystarter/backend"
)

type StartCommand struct {
}

func (c StartCommand) Exec(args ...string) {
	if len(args) > 0 {
		svcName := args[0]
		if svcName == "-all" {
			backend.StartAllServices(args[1:]...)
		} else {
			svc := backend.GetService(svcName, args[1:]...)
			if svc != nil {
				svc.Start()
			}
		}
	} else {
		glg.Error("Specify service name.")
	}
}
