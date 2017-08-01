package commands

import (
	"github.com/vetcher/easystarter/services"
	"github.com/vetcher/easystarter/util"
)

type StopCommand struct {
	allFlag bool
	args    []string
}

func (c *StopCommand) Validate(args ...string) error {
	if len(args) > 0 {
		c.allFlag = util.StrInStrs(ALL, args)
		if c.allFlag {
			c.args = services.ServiceManager.AllServicesNames()
		} else {
			c.args = CompleteNames(args)
		}
		return nil
	}
	return AtLeastOneArgumentErr
}

func (c *StopCommand) Exec() error {
	services.ServiceManager.Stop(c.args...)
	return nil
}
