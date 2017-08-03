package commands

import (
	"github.com/kpango/glg"
	"github.com/vetcher/easystarter/backend/services"
	"github.com/vetcher/easystarter/client/util"
)

type KillCommand struct {
	allFlag bool
	args    []string
}

func (c *KillCommand) Validate(args ...string) error {
	if len(args) > 0 {
		c.allFlag = util.StrInStrs(ALL, args)
		if c.allFlag {
			c.args = <-services.ServeAllServicesNames()
		} else {
			c.args = CompleteNames(args)
		}
		return nil
	}
	return AtLeastOneArgumentErr
}

func (c *KillCommand) Exec() error {
	err := <-services.ServeKillServices(c.args...)
	if err != nil {
		glg.Errorf("Kill error: %v", err)
	}
	return nil
}
