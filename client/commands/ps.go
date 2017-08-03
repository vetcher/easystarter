package commands

import (
	"fmt"
	"strings"
	"time"

	"github.com/gosuri/uitable"
	"github.com/kpango/glg"
	"github.com/vetcher/easystarter/backend/services"
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

func (c *PSCommand) Exec() error {
	glg.Print(printServices(c.allFlag))
	return nil
}

func printServices(allFlag bool) string {
	table := uitable.New()
	table.MaxColWidth = 60
	table.Wrap = true
	table.AddRow("#", glg.White("Service"), "Status", "Command line arguments")
	now := time.Now()
	for i, info := range <-services.ServeServicesInfo(allFlag) {
		upFor := time.Duration(0)
		if !info.StartupTime.IsZero() {
			upFor = now.Sub(info.StartupTime)
		}
		table.AddRow(i, info.Name, fmt.Sprintf("%s %.0fs", info.Status, upFor.Seconds()), strings.Join(info.Args, " "))
	}

	return fmt.Sprintf("In configuration %v services\n%v",
		len(<-services.ServeAllServicesNames()), table.String())
}
