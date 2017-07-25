package backend

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"os"

	"path/filepath"

	"github.com/kpango/glg"
)

var (
	TARGET_PREFIX string = *flag.String("prefix", "", "Use to find service from dir")
	TARGET_SUFFIX string = *flag.String("suffix", "", "Use to find service from dir")
	TARGET_FILE   string = *flag.String("filename", "Makefile", "This file name is used for `go run` command")
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

func RestartAllServices(args ...string) {
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
	println(filepath.Join(TARGET_PREFIX, svcName, TARGET_SUFFIX, TARGET_FILE))
	svc, err := svcByTargetName(filepath.Join(TARGET_PREFIX, svcName, TARGET_SUFFIX, TARGET_FILE))
	if err != nil {
		return nil, err
	}
	svc.Name = svcName
	return svc, nil
}

func findService(svcName string) (*service, error) {
	svc, exist := allServices[svcName]
	if !exist {
		svc, err := svcByNameFromDir(svcName)
		if err != nil {
			return nil, fmt.Errorf("Can't find service %v.", svcName)
		}
		allServices[svcName] = svc
		return svc, nil
	} else {
		return svc, nil
	}
}

func GetService(svcName string, args ...string) *service {
	svc, err := findService(svcName)
	if err != nil {
		glg.Warn(err)
		return nil
	}
	err = svc.SetupService(args...)
	if err != nil {
		glg.Errorf("Can't create service: %v", err)
		return nil
	}
	return svc
}

func StartAllServices(args ...string) {
	for _, svc := range allServices {
		if svc.IsRunning {
			glg.Infof("%v already started.", svc.Name)
		} else {
			svc := GetService(svc.Name, args...)
			if svc != nil {
				go svc.Start()
			}
		}
	}
}

func StopAllServices() {
	for _, svc := range allServices {
		svc.Stop()
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

func KillAllServices() {
	for _, svc := range allServices {
		svc.Kill()
	}
}

func KillService(svcName string) {
	svc, exist := allServices[svcName]
	if !exist {
		glg.Warnf("Can't find service %v.", svcName)
	} else {
		svc.Kill()
	}
}
