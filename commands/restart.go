package commands

import (
	"github.com/vetcher/easystarter/services"
	"github.com/vetcher/easystarter/util"
)

type RestartCommand struct {
	allFlag bool
	args    []string
}

func (c *RestartCommand) Validate(args ...string) error {
	if len(args) > 0 {
		c.allFlag = util.StrInStrs(ALL, args)
		if c.allFlag {
			c.args = services.ServiceManager.AllServicesNames()
		} else {
			c.args = args
		}
		return nil
	}
	return AtLeastOneArgumentErr
}

func (c *RestartCommand) Exec() error {
	services.ServiceManager.Restart(c.args...)
	return nil
}
