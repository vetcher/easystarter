package commands

import (
	"github.com/vetcher/easystarter/backend"
)

type RestartCommand struct {
	allFlag bool
}

func (c RestartCommand) Validate(args ...string) error {
	if len(args) > 0 {
		c.allFlag = args[0] == ALL
		return nil
	}
	return AtLeastOneArgumentErr
}

func (c RestartCommand) Exec(args ...string) error {
	if c.allFlag {
		backend.RestartAllServices(args...)
	} else {
		backend.RestartService(args[0])
	}
	return nil
}
