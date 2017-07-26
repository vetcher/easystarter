package backend

import (
	"errors"
	"flag"
	"fmt"
	"os"

	"path/filepath"

	"github.com/kpango/glg"
)

const (
	CANT_FIND_SERVICE                      = "Can't find service"
	FATAL_WHEN_LOAD_SERVICES_CONFIGURATION = 1
)

var (
	TARGET_PREFIX = flag.String("prefix", "", "Use to find service from dir")
	TARGET_SUFFIX = flag.String("suffix", "", "Use to find service from dir")
	TARGET_FILE   = flag.String("filename", "Makefile", "This file name is used for `go run` command")
)

func RestartAllServices(args ...string) {
	StopAllServices()
	err, fatal := loadServicesConfiguration()
	if err != nil {
		glg.Errorf("Can't load services: %v.", err)
		if fatal {
			os.Exit(FATAL_WHEN_LOAD_SERVICES_CONFIGURATION)
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
	println(filepath.Join(*TARGET_PREFIX, svcName, *TARGET_SUFFIX, *TARGET_FILE))
	svc, err := svcByTargetName(filepath.Join(*TARGET_PREFIX, svcName, *TARGET_SUFFIX, *TARGET_FILE))
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
			return nil, fmt.Errorf("%s %s", CANT_FIND_SERVICE, svcName)
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
			glg.Infof("%v already started", svc.Name)
		} else {
			svc := GetService(svc.Name, args...)
			if svc != nil {
				abs, _ := filepath.Abs("./")
				println("BEFORE START:", abs, svc.Name)
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

func StopAllServicesAndSync() {
	for _, svc := range allServices {
		svc.SyncStop()
	}
}

func StopService(svcName string) {
	svc, exist := allServices[svcName]
	if !exist {
		glg.Warnf("%s %s", CANT_FIND_SERVICE, svcName)
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
		glg.Warnf("%s %s", CANT_FIND_SERVICE, svcName)
	} else {
		svc.Kill()
	}
}
