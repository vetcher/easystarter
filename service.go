package main

import (
	"fmt"
	"os/exec"
)

type service struct {
	Name                  string
	Args                  []string
	CurrentServiceChannel chan interface{}
	CurrentExternalCmd    *exec.Cmd
}

func (svc *service) wait() {
	for {
		select {
		case signal := <-svc.CurrentServiceChannel:
			switch typedSignal := signal.(type) {
			case string:
				if typedSignal == "kill" {
					err := svc.CurrentExternalCmd.Process.Kill()
					if err != nil {
						fmt.Printf("[!] Service %v can't be killed because of %v.\n", svc.Name, err)
					}
					break
				}
			}
		}
	}
}

func (svc *service) Start() {
	go func() {
		err := svc.CurrentExternalCmd
		if err != nil {
			fmt.Printf("[!] Can't start service %v because of %v.\n", svc.Name, err)
			return
		}
		go func() {
			svc.CurrentServiceChannel <- svc.CurrentExternalCmd.Wait()
		}()
		fmt.Printf("[_] Start service %v.\n", svc.Name)
		svc.wait()
	}()
}

var AllServices map[string]*service

func StartNewService(name string, args ...string) {
	service, exist := AllServices[name]
	if exist {
		newServiceChannel := make(chan interface{})
		service.CurrentServiceChannel = newServiceChannel
		args = append(args, service.Args...)
		cmd := exec.Command("go", args...)
		service.CurrentExternalCmd = cmd

	} else {
		fmt.Printf("[!] Can't find service %v.\n", name)
	}
}
