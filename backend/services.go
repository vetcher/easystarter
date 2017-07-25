package backend

import (
	"fmt"
	"strings"
)

func ServicesString(allFlag bool) string {
	var svcStrings []string
	runningCount := 0
	for _, svc := range allServices {
		if svc.IsRunning || allFlag {
			if svc.IsRunning {
				runningCount++
			}
			svcStrings = append(svcStrings, svc.String())
		}
	}

	return fmt.Sprintf("In configuration %v services, %v is up\n%v",
		len(allServices), runningCount, strings.Join(svcStrings, "\n"))
}
