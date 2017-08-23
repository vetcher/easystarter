package commands

import (
	"github.com/kpango/glg"
	"github.com/vetcher/easystarter/backend/services"
	"github.com/vetcher/easystarter/client/util"
)

type StartCommand struct {
	allFlag bool
	args    []string
}

func (c *StartCommand) Validate(args ...string) error {
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

func (c *StartCommand) Exec() error {
	err := <-services.ServeStartServices(c.args...)
	if err != nil {
		glg.Errorf("Start error: %v", err)
	}
	return nil
}
