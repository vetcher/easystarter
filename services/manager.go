package services

import (
	"fmt"
	"sync"
	"time"

	"strings"

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

func (f *serviceManager) Start(svcNames ...string) {
	for _, svcName := range svcNames {
		name := svcName
		go f.start(name)
	}
}

func (f *serviceManager) start(svcName string) {
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
		if svc.IsRunning() {
			glg.Infof("START %s", glg.Yellow(svcName))
		}
	}
}

func (f *serviceManager) Stop(svcNames ...string) {
	var wg sync.WaitGroup
	for _, svcName := range svcNames {
		name := svcName
		wg.Add(1)
		go func() {
			f.stop(name)
			wg.Done()
		}()
	}
	wg.Wait()
}

func (f *serviceManager) stop(svcName string) {
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

func (f *serviceManager) Restart(svcNames ...string) {
	var wg sync.WaitGroup
	for _, svcName := range svcNames {
		name := svcName
		wg.Add(1)
		go func() {
			f.stop(name)
			f.start(name)
			wg.Done()
		}()
	}
	wg.Wait()
}

func (f *serviceManager) Info(allFlag bool) string {
	runningCount := 0
	table := uitable.New()
	table.AddRow("#", glg.White("Service"), "Status", "Command line arguments")
	now := time.Now()
	i := 1
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
			table.AddRow(i, info.Name, fmt.Sprintf("%s %.0fs", info.Status, upFor.Seconds()), strings.Join(info.Args, " "))
			i += 1
		}
	}

	return fmt.Sprintf("In configuration %v services, %v is up\n%v",
		len(f.repo.services), runningCount, table.String())
}

func (f *serviceManager) Kill(svcNames ...string) {
	var wait sync.WaitGroup
	for _, svcName := range svcNames {
		name := svcName
		wait.Add(1)
		go func() {
			f.kill(name)
			wait.Done()
		}()
	}
	wait.Wait()
}

func (f *serviceManager) kill(svcName string) {
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

func (f *serviceManager) AllServicesNames() []string {
	return f.repo.Names()
}
