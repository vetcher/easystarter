package commands

import (
	"github.com/kpango/glg"
	"github.com/vetcher/easystarter/services"
)

type PSCommand struct {
	allFlag bool
}

func (c *PSCommand) Validate(args ...string) error {
	c.allFlag = false
	if len(args) > 0 {
		c.allFlag = args[0] == ALL
	}
	return nil
}

func (c *PSCommand) Exec(args ...string) error {
	glg.Print(services.ServiceManager.Info(c.allFlag))
	return nil
}
