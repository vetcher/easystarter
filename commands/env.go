package commands

import (
	"github.com/kpango/glg"
	"github.com/vetcher/easystarter/backend"
)

type EnvCommand struct {
}

func (c EnvCommand) Exec(args ...string) {
	if len(args) > 0 {
		if args[0] == "-all" {
			glg.Info(backend.CurrentEnvironmentString())
		} else if args[0] == "-reload" {
			if backend.SetupEnv() {
				glg.Info("Environment was reloaded.")
			}
		} else {
			glg.Info(backend.AllEnvironmentString())
		}
	}
}
