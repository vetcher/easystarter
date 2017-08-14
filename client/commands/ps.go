package commands

import (
	"fmt"
	"time"

	"github.com/gosuri/uitable"
	"github.com/kpango/glg"
	"github.com/vetcher/easystarter/backend/services"
	"sort"
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

type serviceInfoSlice struct {
	data []services.ServiceInfo
}

func (s *serviceInfoSlice) Len() int {
	return len(s.data)
}

func (s *serviceInfoSlice) Less(i, j int) bool {
	return s.data[i].Name < s.data[j].Name
}

func (s *serviceInfoSlice) Swap(i, j int) {
	s.data[i], s.data[j] = s.data[j], s.data[i]
}

func printServices(allFlag bool) string {
	slice := serviceInfoSlice{<-services.ServeServicesInfo(allFlag)}
	sort.Sort(&slice)
	table := uitable.New()
	table.MaxColWidth = 60
	table.Wrap = true
	table.AddRow("#", glg.White("Service"), "Status")
	now := time.Now()
	for i, info := range slice.data {
		upFor := time.Duration(0)
		if !info.StartupTime.IsZero() {
			upFor = now.Sub(info.StartupTime)
		}
		table.AddRow(i+1, info.Name, fmt.Sprintf("%s %.0fs", info.Status, upFor.Seconds()))
	}

	return fmt.Sprintf("In configuration %v services\n%v",
		len(<-services.ServeAllServicesNames()), table.String())
}
