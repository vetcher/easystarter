package services

import (
	"github.com/vetcher/easystarter/backend"
)

func ServeServicesInfo(allFlag bool) <-chan []ServiceInfo {
	info := serviceManager.Info(allFlag)
	respChan := make(chan []ServiceInfo, 1)
	go func() {
		respChan <- info
	}()
	return respChan
}

func ServeReloadEnv() <-chan error {
	respChan := make(chan error, 1)
	err := backend.ReloadEnvironment()
	go func() {
		respChan <- err
	}()
	return respChan
}

func ServeGetEnv(allFlag bool) <-chan []string {
	respChan := make(chan []string, 1)
	var env []string
	if allFlag {
		env, _ = backend.AllEnvironmentString()
	} else {
		env = backend.CurrentEnvironmentString()
	}
	go func() {
		respChan <- env
	}()
	return respChan
}

func ServeStartServices(svcNames ...string) <-chan error {
	err := serviceManager.Start(svcNames...)
	respChan := make(chan error, 1)
	go func() {
		respChan <- err
	}()
	return respChan
}

func ServeStopServices(svcNames ...string) <-chan error {
	err := serviceManager.Stop(svcNames...)
	respChan := make(chan error, 1)
	go func() {
		respChan <- err
	}()
	return respChan
}

func ServeKillServices(svcNames ...string) <-chan error {
	err := serviceManager.Kill(svcNames...)
	respChan := make(chan error, 1)
	go func() {
		respChan <- err
	}()
	return respChan
}

func ServeRestartServices(svcNames ...string) <-chan error {
	err := serviceManager.Restart(svcNames...)
	respChan := make(chan error, 1)
	go func() {
		respChan <- err
	}()
	return respChan
}

func ServeLoadServices() <-chan error {
	err := loadServices()
	respChan := make(chan error, 1)
	go func() {
		respChan <- err
	}()
	return respChan
}

func ServeAllServicesNames() <-chan []string {
	names := serviceManager.AllServicesNames()
	respChan := make(chan []string, 1)
	go func() {
		respChan <- names
	}()
	return respChan
}
