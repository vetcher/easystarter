package services

import (
	"fmt"

	"sync"

	"strings"

	"github.com/kpango/glg"
)

const CANT_FIND_SERVICE = "Can't find service"

var (
	ServiceManager *serviceManager
)

type serviceManager struct {
	repo *ServiceRepository
}

func init() {
	ServiceManager = &serviceManager{
		repo: NewRepository(),
	}
}

func (f *serviceManager) RegisterService(config *ServiceConfig) error {
	return f.repo.RegisterService(config)
}

func (f *serviceManager) StartAllServices() {
	for _, svcName := range f.repo.Names() {
		go f.Start(svcName)
	}
}

func (f *serviceManager) Start(svcName string) {
	svc, err := f.repo.GetService(svcName)
	if err != nil {
		glg.Errorf("Start %s error: %v", svcName, err)
	}
	errc := make(chan error)
	defer close(errc)

	go func() {
		err := svc.Build()
		if err != nil {
			errc <- fmt.Errorf("build error: %v", err)
			return
		}
		err = svc.Start()
		if err != nil {
			errc <- fmt.Errorf("start error: %v", err)
			return
		}
		errc <- nil
		return
	}()

	err = <-errc
	if err != nil {
		glg.Errorf("Start %s error: %v", svcName, err)
	} else {
		glg.Infof("START %s", glg.Yellow(svcName))
	}
}

func (f *serviceManager) StopAllServices() {
	var wait sync.WaitGroup
	for _, svcName := range f.repo.Names() {
		wait.Add(1)
		go func() {
			f.Stop(svcName)
			wait.Done()
		}()
	}
	wait.Wait()
}

func (f *serviceManager) Stop(svcName string) {
	svc, err := f.repo.GetService(svcName)
	if err != nil {
		glg.Warnf("Stop %s error: %v", svcName, err)
		return
	}
	err = svc.Stop()
	if err != nil {
		glg.Errorf("Stop %s error: %v", svcName, err)
		return
	}
	svc.Sync()
	glg.Infof("STOP %v", glg.Yellow(svcName))
}

func (f *serviceManager) Restart(svcName string) {
	f.Stop(svcName)
	f.Start(svcName)
}

func (f *serviceManager) RestartAllServices() {
	var wait sync.WaitGroup
	for _, svcName := range f.repo.Names() {
		wait.Add(1)
		go func() {
			f.Restart(svcName)
			wait.Done()
		}()
	}
	wait.Wait()
}

func (f *serviceManager) String(allFlag bool) string {
	var svcStrings []string
	runningCount := 0
	for _, svc := range f.repo.services {
		if svc.IsRunning() || allFlag {
			if svc.IsRunning() {
				runningCount++
			}
			svcStrings = append(svcStrings, svc.String())
		}
	}

	return fmt.Sprintf("In configuration %v services, %v is up\n%v",
		len(f.repo.services), runningCount, strings.Join(svcStrings, "\n"))
}
