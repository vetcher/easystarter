package services

import (
	"fmt"
	"os"
	"os/exec"
	"sync"
	"time"

	"github.com/vetcher/easystarter/printer"
)

type service struct {
	// Имя сервиса
	Name string `json:"-"`
	// Аргументы, которые будут переданы в сервис как параметры командной строки
	Args []string `json:"args"`
	// Канал, принимающий сообщения для действий
	currentServiceChannel chan interface{} `json:"-"`
	currentExternalCmd    *exec.Cmd        `json:"-"`
	// Отметка о старте
	StartTime time.Time `json:"-"`
	IsRunning bool      `json:"-"`
	// *.go файл, который запускает сервис
	Target string `json:"target"`
	// Не даёт запустить новый сервис, если старый еще не закончил работу
	syncMutex sync.Mutex `json:"-"`
}

var allServices map[string]*service

// Returns true if process should be stopped
func (svc *service) stringSwitch(text string) bool {
	switch text {
	case OK_SIGNAL:
		return true
	case KILL_SIGNAL:
		err := svc.currentExternalCmd.Process.Kill()
		if err != nil {
			printer.Printf("!", "Service %v can't be killed because %v.", svc.Name, err)
		}
		return true
	}
	return false
}

func (svc *service) wait() {
	svc.IsRunning = true
	for {
		select {
		case signal := <-svc.currentServiceChannel:
			switch typedSignal := signal.(type) {
			case string:
				if svc.stringSwitch(typedSignal) {
					return
				}
			case error:
				printer.Printf("!", "Error with service %v:", svc.Name)
				printer.Printf("!", "%v", typedSignal)
				return
			}
		}
	}
}

func (svc *service) Start() {
	go func() {
		svc.StartTime = time.Now()
		out, err := os.Create(fmt.Sprintf("logs/%v.log", svc.Name))
		// Init log file and all output would write to file
		// If init unsuccessful out will be written to Stdout and Stderr
		if err != nil {
			printer.Printf("?", "Can't init %v.log file because %v", svc.Name, err)
			svc.currentExternalCmd.Stdout = os.Stdout
			svc.currentExternalCmd.Stderr = os.Stderr
		} else {
			svc.currentExternalCmd.Stdout = out
			svc.currentExternalCmd.Stderr = out
		}
		err = svc.currentExternalCmd.Start()
		if err != nil {
			printer.Print("!", "Can't start service %v because %v.", svc.Name, err)
			return
		}
		svc.syncMutex.Lock()
		go func() {
			err := svc.currentExternalCmd.Wait()
			if err != nil {
				svc.currentServiceChannel <- err
			} else {
				svc.currentServiceChannel <- OK_SIGNAL
			}
		}()
		printer.Printf("I", "Start %v.", svc.Name)
		svc.wait()
		svc.cleanService()
		svc.StartTime = time.Time{}
		printer.Printf("I", "Stop %v.", svc.Name)
	}()
}

func (svc *service) cleanService() {
	close(svc.currentServiceChannel)
	svc.currentServiceChannel = nil
	svc.currentExternalCmd = nil
	svc.IsRunning = false
	svc.syncMutex.Unlock()
}

func (svc *service) buildService() error {
	buildCmd := exec.Command("go", "build", "-o", "./bin/"+svc.Name, svc.Target)
	err := buildCmd.Start()
	if err != nil {
		return fmt.Errorf("can't build because %v", err)
	}
	err = buildCmd.Wait()
	if err != nil {
		return fmt.Errorf("can't build because %v", err)
	}
	return nil
}

func (svc *service) SetupService(args ...string) error {
	if svc.currentExternalCmd != nil || svc.currentServiceChannel != nil {
		return fmt.Errorf("service %v already in use", svc.Name)
	} else {
		err := svc.buildService()
		if err != nil {
			return err
		}
		svc.currentServiceChannel = make(chan interface{})
		svc.Args = append(svc.Args, args...)
		runArgs := []string{}
		runArgs = append(runArgs, svc.Args...)
		cmd := exec.Command("./bin/"+svc.Name, runArgs...)
		svc.currentExternalCmd = cmd
		return nil
	}
}
