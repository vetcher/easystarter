package commands

import (
	"github.com/vetcher/easystarter/backend/services"
	"github.com/vetcher/easystarter/client/util"
	"github.com/kpango/glg"
)

type StopCommand struct {
	allFlag bool
	args    []string
}

func (c *StopCommand) Validate(args ...string) error {
	c.allFlag = false
	c.args = []string{}
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

func (c *StopCommand) Exec() error {
	err := <-services.ServeStopServices(c.args...)
	if err != nil {
		glg.Errorf("Stop error: %v", err)
	}
	return nil
}
