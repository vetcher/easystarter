package commands

import (
	"github.com/vetcher/easystarter/services"
)

type StartCommand struct {
	allFlag bool
}

func (c *StartCommand) Validate(args ...string) error {
	if len(args) > 0 {
		c.allFlag = args[0] == ALL
		return nil
	}
	return AtLeastOneArgumentErr
}

func (c *StartCommand) Exec(args ...string) error {
	if c.allFlag {
		services.ServiceManager.StartAllServices()
	} else {
		services.ServiceManager.Start(args[0])
	}
	return nil
}
