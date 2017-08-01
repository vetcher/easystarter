package commands

import (
	"github.com/vetcher/easystarter/services"
	"github.com/vetcher/easystarter/util"
)

type StartCommand struct {
	allFlag bool
	args    []string
}

func (c *StartCommand) Validate(args ...string) error {
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

func (c *StartCommand) Exec() error {
	services.ServiceManager.Start(c.args...)
	return nil
}
