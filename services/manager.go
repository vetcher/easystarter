package services

import (
	"fmt"
	"sync"

	"time"

	"github.com/gosuri/uitable"
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
	err := loadServices()
	if err != nil {
		glg.Fatalf("can't load services: %v", err)
	}
}

func (f *serviceManager) RegisterService(config *ServiceConfig) error {
	return f.repo.RegisterService(config)
}

func (f *serviceManager) StartAllServices() {
	for _, svcName := range f.repo.Names() {
		name := svcName
		go f.Start(name)
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
		name := svcName
		wait.Add(1)
		go func() {
			f.Stop(name)
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
	if svc.IsRunning() {
		err = svc.Stop()
		if err != nil {
			glg.Errorf("Stop %s error: %v", svcName, err)
			return
		}
		svc.Sync()
		glg.Infof("STOP %v", glg.Yellow(svcName))
	}
}

func (f *serviceManager) Restart(svcName string) {
	f.Stop(svcName)
	f.Start(svcName)
}

func (f *serviceManager) RestartAllServices() {
	var wait sync.WaitGroup
	for _, svcName := range f.repo.Names() {
		name := svcName
		wait.Add(1)
		go func() {
			f.Restart(name)
			wait.Done()
		}()
	}
	wait.Wait()
}

func (f *serviceManager) Info(allFlag bool) string {
	runningCount := 0
	table := uitable.New()
	now := time.Now()
	for _, svc := range f.repo.services {
		if svc.IsRunning() || allFlag {
			if svc.IsRunning() {
				runningCount++
			}
			info := svc.Info()
			upFor := time.Duration(0)
			if !info.StartupTime.IsZero() {
				upFor = now.Sub(info.StartupTime)
			}
			table.AddRow(info.Name, info.Status, upFor, info.Args)
		}
	}

	return fmt.Sprintf("In configuration %v services, %v is up\n%v",
		len(f.repo.services), runningCount, table)
}

func (f *serviceManager) KillAllServices() {
	var wait sync.WaitGroup
	for _, svcName := range f.repo.Names() {
		name := svcName
		wait.Add(1)
		go func() {
			f.Kill(name)
			wait.Done()
		}()
	}
	wait.Wait()
}

func (f *serviceManager) Kill(svcName string) {
	svc, err := f.repo.GetService(svcName)
	if err != nil {
		glg.Warnf("Kill %s error: %v", svcName, err)
		return
	}
	err = svc.Kill()
	if err != nil {
		glg.Errorf("Kill %s error: %v", svcName, err)
		return
	}
	svc.Sync()
	glg.Infof("KILL %v", glg.Yellow(svcName))
}
