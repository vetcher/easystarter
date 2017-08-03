package commands

import (
	"github.com/kpango/glg"
	"github.com/vetcher/easystarter/backend/services"
	"github.com/vetcher/easystarter/client/util"
)

type RestartCommand struct {
	allFlag bool
	args    []string
}

func (c *RestartCommand) Validate(args ...string) error {
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

func (c *RestartCommand) Exec() error {
	err := <-services.ServeRestartServices(c.args...)
	if err != nil {
		glg.Errorf("Restart error: %v", err)
	}
	return nil
}
