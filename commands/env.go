package commands

import (
	"github.com/kpango/glg"
	"github.com/vetcher/easystarter/backend"
	"errors"
)

type EnvCommand struct {
}

func (c EnvCommand) Exec(args ...string) error {
	if len(args) > 0 {
		if args[0] == "-all" {
			glg.Info(backend.CurrentEnvironmentString())
		} else if args[0] == "-reload" {
			if backend.SetupEnv() {
				glg.Info("Environment was reloaded.")
			} else {
				return errors.New("Can't load environment")
			}
		} else {
			glg.Info(backend.AllEnvironmentString())
		}
	}
	return nil
}
