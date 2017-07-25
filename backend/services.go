package backend

import (
	"fmt"
	"strings"
	"time"
)

func ServicesString(param string) string {
	isPrintAll := false
	if param == "-all" {
		isPrintAll = true
	}
	now := time.Now()
	var svcStrings []string
	runningCount := 0
	for key, val := range allServices {
		if val.IsRunning || isPrintAll {
			isRunningStr := "Down"
			if val.IsRunning {
				isRunningStr = fmt.Sprintf("Up for %v", now.Sub(val.StartTime))
				runningCount++
			}
			svcStrings = append(svcStrings, fmt.Sprintf("%v %v %v", key, val.Args, isRunningStr))
		}
	}

	return fmt.Sprintf("In configuration %v services, %v is up\n%v",
		len(allServices), runningCount, strings.Join(svcStrings, "\n"))
}
