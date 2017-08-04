package commands

import (
	"fmt"

	"github.com/gosuri/uitable"
	"github.com/kpango/glg"
	"github.com/vetcher/easystarter/backend/services"
	"github.com/vetcher/easystarter/client/util"
)

type CfgCommand struct {
	allFlag    bool
	reloadFlag bool
	svcNames   []string
}

func (c *CfgCommand) Validate(args ...string) error {
	if len(args) > 0 {
		c.reloadFlag = util.StrInStrs(RELOAD, args)
		c.allFlag = util.StrInStrs(ALL, args)
		switch {
		case c.allFlag:
			return nil
		case c.reloadFlag:
			pos := util.IndexStrInStrs(RELOAD, args)
			args[pos] = args[len(args)-1]
			c.svcNames = args[:len(args)-1]
		}
		if c.allFlag {
			return nil
		}

	}
	return nil
}

func (c *CfgCommand) Exec() error {
	if c.reloadFlag {
		if c.allFlag {
			err := <-services.ServeLoadAllServices()
			if err != nil {
				glg.Errorf("Reload config error: %v", err)
			}
		} else {
			err := <-services.ServeLoadServices(c.svcNames...)
			if err != nil {
				glg.Errorf("Reload config error: %v", err)
			}
		}
		glg.Success("Reload config: OK")
	}
	glg.Printf("\n%s", c.configuration())
	return nil
}

func (c *CfgCommand) configuration() string {
	table := uitable.New()
	table.MaxColWidth = 100
	table.Wrap = true
	infos := <-services.ServeServicesInfo(true)
	for _, info := range infos {
		table.AddRow("Service:", fmt.Sprintf("%s:%s", info.Name, info.Version))
		table.AddRow("Dir:", util.StringOrEmpty(info.Dir))
		table.AddRow("Target:", util.StringOrEmpty(info.Target))
		table.AddRow("Args:", info.Args)
		table.AddRow("")
	}
	return table.String()
}
