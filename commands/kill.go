package commands

import (
	"github.com/vetcher/easystarter/services"
	"github.com/vetcher/easystarter/util"
)

type KillCommand struct {
	allFlag bool
	args    []string
}

func (c *KillCommand) Validate(args ...string) error {
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

func (c *KillCommand) Exec() error {
	services.ServiceManager.Kill(c.args...)
	return nil
}
