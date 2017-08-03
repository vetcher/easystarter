package services

import (
	"fmt"

	"github.com/kpango/glg"
	"github.com/vetcher/easystarter/backend/util"
)

const CANT_FIND_SERVICE = "Can't find service"

type Manager interface {
	RegisterService(config *ServiceConfig) error // Add new service to configuration.
	Start(svcNames ...string) error              // Start services with specified names.
	Stop(svcNames ...string) error               // Normal and safely stop services with specified names.
	Restart(svcNames ...string) error            // Restart services with specified names.
	Kill(svcNames ...string) error               // Fast and rude stop services with specified names.
	Info(allFlag bool) []ServiceInfo             // Information about configured services.
	AllServicesNames() []string                  // List of all services in configuration.
}

var (
	serviceManager Manager
)

type manager struct {
	repo         *ServiceRepository
	startSteps   []Step
	stopSteps    []Step
	killSteps    []Step
	restartSteps []Step
}

func init() {
	serviceManager = &manager{
		repo: NewRepository(),
		startSteps: []Step{
			//{
			//	Name: "checkout",
			//	Do:   SwitchVersion,
			//},
			//{
			//	Name: "make",
			//	Do:   CallMakeDepGen,
			//},
			{
				Name: "prepare",
				Do:   Service.Prepare,
			},
			{
				Name: "build",
				Do:   Service.Build,
			},
			{
				Name: "start",
				Do:   Service.Start,
			},
		},
		stopSteps: []Step{
			{
				Name: "stop",
				Do:   Service.Stop,
			},
			{
				Name: "sync",
				Do:   Service.Sync,
			},
		},
		killSteps: []Step{
			{
				Name: "kill",
				Do:   Service.Kill,
			},
			{
				Name: "sync",
				Do:   Service.Sync,
			},
		},
		restartSteps: []Step{
			{
				Name: "stop",
				Do:   Service.Stop,
			},
			{
				Name: "sync",
				Do:   Service.Sync,
			},
			{
				Name: "prepare",
				Do:   Service.Prepare,
			},
			{
				Name: "build",
				Do:   Service.Build,
			},
			{
				Name: "start",
				Do:   Service.Start,
			},
		},
	}
	err := loadServices()
	if err != nil {
		glg.Fatalf("can't load services: %v", err)
	}
}

func (f *manager) RegisterService(config *ServiceConfig) error {
	return f.repo.RegisterService(config)
}

func (f *manager) Start(svcNames ...string) error {
	errc := make(chan error, len(svcNames))
	defer close(errc)
	for _, svcName := range svcNames {
		name := svcName
		go func() {
			errc <- f.start(name)
		}()
	}
	var multiErr []error
	for range svcNames {
		err, ok := <-errc
		if ok {
			multiErr = append(multiErr, err)
		} else {
			break
		}
	}
	return util.ComposeErrors(multiErr)
}

func (f *manager) start(svcName string) error {
	svc, err := f.repo.GetService(svcName)
	if err != nil {
		return fmt.Errorf("%s start error: %v", svcName, err)
	}
	err = CallStepByStep(svc, f.startSteps...)
	if err != nil {
		return fmt.Errorf("%s: %v", svcName, err)
	}
	if !svc.IsRunning() {
		return fmt.Errorf("%s: service not started", svcName)
	}
	return nil
}

func (f *manager) Stop(svcNames ...string) error {
	errc := make(chan error, len(svcNames))
	defer close(errc)
	for _, svcName := range svcNames {
		name := svcName
		go func() {
			errc <- f.stop(name)
		}()
	}
	var multiErr []error
	for range svcNames {
		err, ok := <-errc
		if ok {
			multiErr = append(multiErr, err)
		} else {
			break
		}
	}
	return util.ComposeErrors(multiErr)
}

func (f *manager) stop(svcName string) error {
	svc, err := f.repo.GetService(svcName)
	if err != nil {
		return fmt.Errorf("%s: stop error: %v", svcName, err)
	}
	if svc.IsRunning() {
		err = CallStepByStep(svc, f.stopSteps...)
		if err != nil {
			return fmt.Errorf("%s: %v", svcName, err)
		}
	}
	return nil
}

func (f *manager) Restart(svcNames ...string) error {
	errc := make(chan error, len(svcNames))
	defer close(errc)
	for _, svcName := range svcNames {
		name := svcName
		go func() {
			errc <- f.restart(name)
		}()
	}
	var multiErr []error
	for range svcNames {
		err, ok := <-errc
		if ok {
			multiErr = append(multiErr, err)
		} else {
			break
		}
	}
	return util.ComposeErrors(multiErr)
}

func (f *manager) restart(svcName string) error {
	svc, err := f.repo.GetService(svcName)
	if err != nil {
		return fmt.Errorf("%s start error: %v", svcName, err)
	}
	err = CallStepByStep(svc, f.restartSteps...)
	if err != nil {
		return fmt.Errorf("%s: %v", svcName, err)
	}
	if !svc.IsRunning() {
		return fmt.Errorf("%s: service not started", svcName)
	}
	return nil
}

func (f *manager) Kill(svcNames ...string) error {
	errc := make(chan error, len(svcNames))
	defer close(errc)
	for _, svcName := range svcNames {
		name := svcName
		go func() {
			errc <- f.kill(name)
		}()
	}
	var multiErr []error
	for range svcNames {
		err, ok := <-errc
		if ok {
			multiErr = append(multiErr, err)
		} else {
			break
		}
	}
	return util.ComposeErrors(multiErr)
}

func (f *manager) Info(allFlag bool) []ServiceInfo {
	var infos []ServiceInfo
	for _, svc := range f.repo.services {
		if svc.IsRunning() || allFlag {
			infos = append(infos, *svc.Info())
		}
	}
	return infos
}

func (f *manager) kill(svcName string) error {
	svc, err := f.repo.GetService(svcName)
	if err != nil {
		return fmt.Errorf("%s: kill error: %v", svcName, err)
	}
	err = CallStepByStep(svc, f.killSteps...)
	if err != nil {
		return fmt.Errorf("%s: %v", svcName, err)
	}
	return nil
}

func (f *manager) AllServicesNames() []string {
	return f.repo.Names()
}
