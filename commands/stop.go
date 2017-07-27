package commands

import (
	"github.com/vetcher/easystarter/services"
)

type StopCommand struct {
	allFlag bool
}

func (c *StopCommand) Validate(args ...string) error {
	if len(args) > 0 {
		c.allFlag = args[0] == ALL
		return nil
	}
	return AtLeastOneArgumentErr
}

func (c *StopCommand) Exec(args ...string) error {
	if c.allFlag {
		services.ServiceManager.StopAllServices()
	} else {
		services.ServiceManager.Stop(args[0])
	}
	return nil
}
