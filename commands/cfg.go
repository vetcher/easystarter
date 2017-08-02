package commands

import (
	"github.com/kpango/glg"
	"github.com/vetcher/easystarter/services"
	"github.com/vetcher/easystarter/util"
)

type CfgCommand struct {
	allFlag    bool
	reloadFlag bool
}

func (c *CfgCommand) Validate(args ...string) error {
	if len(args) > 0 {
		c.reloadFlag = util.StrInStrs(RELOAD, args)
		return nil
	}
	return nil
}

func (c *CfgCommand) Exec() error {
	if c.reloadFlag {
		err := services.LoadServices()
		if err != nil {
			glg.Errorf("Reload services: %v", err)
		}
		glg.Success("Reloading config")
	}
	glg.Printf("\n%s", services.ServiceManager.Configuration())
	return nil
}
