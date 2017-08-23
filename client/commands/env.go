package commands

import (
	"fmt"

	"github.com/kpango/glg"
	"github.com/vetcher/easystarter/backend/services"
	"github.com/vetcher/easystarter/client/util"
)

type EnvCommand struct {
	allFlag    bool
	reloadFlag bool
}

func (c *EnvCommand) Validate(args ...string) error {
	if len(args) > 0 {
		c.allFlag = util.StrInStrs(ALL, args)
		c.reloadFlag = util.StrInStrs(RELOAD, args)
		return nil
	}
	return nil
}

func (c *EnvCommand) Exec() error {
	if c.reloadFlag {
		err := <-services.ServeReloadEnv()
		if err != nil {
			return fmt.Errorf("Can't load environment: %v", err)
		}
		glg.Info("Environment was reloaded.")
	}
	glg.Infof("%v", <-services.ServeGetEnv(c.allFlag))
	return nil
}
