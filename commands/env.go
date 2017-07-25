package commands

import (
	"errors"

	"github.com/kpango/glg"
	"github.com/vetcher/easystarter/backend"
)

type EnvCommand struct {
	allFlag    bool
	reloadFlag bool
}

func (c *EnvCommand) Validate(args ...string) error {
	if len(args) > 0 {
		c.allFlag = args[0] == ALL
		c.reloadFlag = args[0] == RELOAD
		return nil
	}
	return AtLeastOneArgumentErr
}

func (c *EnvCommand) Exec(args ...string) error {
	if c.allFlag {
		glg.Info(backend.AllEnvironmentString())
	} else if c.reloadFlag {
		if backend.SetupEnv() {
			glg.Info("Environment was reloaded.")
		} else {
			return errors.New("Can't load environment")
		}
	} else {
		glg.Info(backend.CurrentEnvironmentString())
	}
	return nil
}
