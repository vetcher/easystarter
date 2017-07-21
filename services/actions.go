package services

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"time"

	"github.com/vetcher/easystarter/printer"
)

const (
	OK_SIGNAL   = "ok"
	KILL_SIGNAL = "kill"
)

var (
	TARGET_PREFIX string = *flag.String("prefix", "", "Use to find service from dir")
	TARGET_SUFFIX string = *flag.String("suffix", "cmd/", "Use to find service from dir")
	TARGET_FILE   string = *flag.String("filename", "main.go", "This file name is used for `go run` command")
)

func loadServices() (error, bool) {
	allServices = make(map[string]*service)
	raw, err := ioutil.ReadFile("services.json")
	if err != nil {
		return err, false
	}
	err = json.Unmarshal(raw, &allServices)
	if err != nil {
		return err, false
	}
	for key, val := range allServices {
		if key == "-all" || key == "-env" {
			return fmt.Errorf("name `%v` is not allowed", key), true
		}
		if val.Target == "" {
			printer.Printf("?", "Field `target` is not provided for %v service", key)
			delete(allServices, key)
		} else {
			val.Name = key
		}
	}
	return nil, false
}

func init() {
	err, fatal := loadServices()
	if err != nil {
		printer.Printf("!", "Can't load services because %v.", err)
		if fatal {
			os.Exit(1)
		}
	}
}

func ReloadServices(args ...string) {
	StopAll()
	err, fatal := loadServices()
	if err != nil {
		printer.Printf("!", "Can't load services because %v.", err)
		if fatal {
			os.Exit(1)
		}
	}
	StartAll(args...)
}

func ReloadService(svcName string) {
	StopService(svcName)
	svc := NewService(svcName)
	if svc != nil {
		svc.Start()
	}
}

func isTargetExist(target string) bool {
	fileInfo, err := os.Stat(target)
	if os.IsNotExist(err) || fileInfo.IsDir() {
		return false
	}
	return true
}

func svcByTargetName(target string) (*service, error) {
	if !isTargetExist(target) {
		return nil, errors.New("service does not exist")
	}
	return &service{
		Target: target,
	}, nil
}

// Ищет сервис в соседних директориях
// Запускаемый файл должен называться `main.go`
func svcByNameFromDir(svcName string) (*service, error) {
	svc, err := svcByTargetName(fmt.Sprintf("%v%v/%v%v/%v", TARGET_PREFIX, svcName, TARGET_SUFFIX, svcName, TARGET_FILE))
	if err != nil {
		return nil, err
	}
	svc.Name = svcName
	return svc, nil
}

func findService(svcName string) *service {
	svc, exist := allServices[svcName]
	if !exist {
		svc, err := svcByNameFromDir(svcName)
		if err != nil {
			printer.Printf("!", "Can't find service %v.", svcName)
			return nil
		}
		allServices[svcName] = svc
		return svc
	} else {
		return svc
	}
}

func NewService(svcName string, args ...string) *service {
	svc := findService(svcName)
	if svc == nil {
		return nil
	}
	err := svc.SetupService(args...)
	if err != nil {
		printer.Printf("?", "Can't create service because %v", err)
		return nil
	}
	return svc
}

func ListServices() string {
	now := time.Now()
	var svcStrs []string
	runningCount := 0
	for key, val := range allServices {
		isRunningStr := "Down"
		if val.IsRunning {
			isRunningStr = fmt.Sprintf("Up for %v", now.Sub(val.StartTime))
			runningCount++
		}
		svcStrs = append(svcStrs, fmt.Sprintf("%v %v %v", key, val.Args, isRunningStr))
	}

	return fmt.Sprintf("In configuration %v services, %v is up\n%v",
		len(allServices), runningCount, strings.Join(svcStrs, "\n"))
}

func StartAll(args ...string) {
	for key, val := range allServices {
		if val.IsRunning {
			printer.Print("?", "%v already started.", key)
		} else {
			svc := NewService(key, args...)
			if svc != nil {
				svc.Start()
			}
		}
	}
}

func StopAll() {
	for key, val := range allServices {
		if val.IsRunning {
			StopService(key)
		}
	}
}

func StopService(svcName string) {
	svc, exist := allServices[svcName]
	if !exist {
		printer.Printf("!", "Can't find service %v.", svcName)
	} else {
		if svc.IsRunning {
			svc.currentServiceChannel <- KILL_SIGNAL
			svc.syncMutex.Lock()
			svc.syncMutex.Unlock()
		}
	}
}
