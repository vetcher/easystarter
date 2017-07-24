package backend

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/kpango/glg"
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

func loadServicesConfiguration() (error, bool) {
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
		if key == "-all" {
			return fmt.Errorf("name `%v` is not allowed", key), true
		}
		if val.Target == "" {
			glg.Warnf("Field `target` is not provided for %v service", key)
			delete(allServices, key)
		} else {
			val.Name = key
		}
	}
	return nil, false
}

func init() {
	err, fatal := loadServicesConfiguration()
	if err != nil {
		glg.Errorf("Can't load services: %v.", err)
		if fatal {
			os.Exit(1)
		}
	}
}

func RestartServices(args ...string) {
	StopAllServices()
	err, fatal := loadServicesConfiguration()
	if err != nil {
		glg.Errorf("Can't load services: %v.", err)
		if fatal {
			os.Exit(1)
		}
	}
	StartAllServices(args...)
}

func RestartService(svcName string) {
	StopService(svcName)
	svc := GetService(svcName)
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
			glg.Warnf("Can't find service %v.", svcName)
			return nil
		}
		allServices[svcName] = svc
		return svc
	} else {
		return svc
	}
}

func GetService(svcName string, args ...string) *service {
	svc := findService(svcName)
	if svc == nil {
		return nil
	}
	err := svc.SetupService(args...)
	if err != nil {
		glg.Errorf("Can't create service: %v", err)
		return nil
	}
	return svc
}

func StartAllServices(args ...string) {
	for key, val := range allServices {
		if val.IsRunning {
			glg.Infof("%v already started.", key)
		} else {
			svc := GetService(key, args...)
			if svc != nil {
				svc.Start()
			}
		}
	}
}

func StopAllServices() {
	for key, val := range allServices {
		if val.IsRunning {
			StopService(key)
		}
	}
}

func StopService(svcName string) {
	svc, exist := allServices[svcName]
	if !exist {
		glg.Warnf("Can't find service %v.", svcName)
	} else {
		svc.Stop()
	}
}
