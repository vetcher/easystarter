package commands

import "github.com/vetcher/easystarter/services"

type KillCommand struct {
	allFlag bool
}

func (c *KillCommand) Validate(args ...string) error {
	if len(args) > 0 {
		c.allFlag = args[0] == ALL
		return nil
	}
	return AtLeastOneArgumentErr
}

func (c *KillCommand) Exec(args ...string) error {
	if c.allFlag {
		services.ServiceManager.KillAllServices()
	} else {
		services.ServiceManager.Kill(args[0])
	}
	return nil
}
