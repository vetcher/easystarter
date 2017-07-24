package backend

import (
	"fmt"
	"os"
	"os/exec"
	"sync"
	"time"

	"github.com/kpango/glg"
	"path/filepath"
)

type service struct {
	// Имя сервиса
	Name string
	// Аргументы, которые будут переданы в сервис как параметры командной строки
	Args []string
	// Канал, принимающий сообщения для действий
	currentServiceChannel chan interface{}
	currentExternalCmd    *exec.Cmd
	// Отметка о старте
	StartTime time.Time
	Target    string
	IsRunning bool
	// Не даёт запустить новый сервис, если старый еще не закончил работу
	syncMutex sync.Mutex
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
			glg.Errorf("Service %v can't be killed: %v.", svc.Name, err)
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
				glg.Errorf("Error with service %v:", svc.Name)
				glg.Errorf("%v", typedSignal)
				return
			}
		}
	}
}

func (svc *service) Start() {
	go func() {
		out, err := os.Create(fmt.Sprintf("./logs/%v.log", svc.Name))
		// Init log file and all output would write to file
		// If init unsuccessful out will be written to Stdout and Stderr
		if err != nil {
			glg.Warnf("Can't init %v.log file: %v", svc.Name, err)
			svc.currentExternalCmd.Stdout = os.Stdout
			svc.currentExternalCmd.Stderr = os.Stderr
		} else {
			svc.currentExternalCmd.Stdout = out
			svc.currentExternalCmd.Stderr = out
		}
		err = svc.currentExternalCmd.Start()
		if err != nil {
			glg.Errorf("Can't start service %v: %v.", svc.Name, err)
			return
		}
		svc.StartTime = time.Now()
		svc.syncMutex.Lock()
		go func() {
			err := svc.currentExternalCmd.Wait()
			if err != nil {
				svc.currentServiceChannel <- err
			} else {
				svc.currentServiceChannel <- OK_SIGNAL
			}
		}()
		glg.Infof("Start %v.", svc.Name)
		svc.wait()
		svc.cleanService()
		svc.StartTime = time.Time{}
		glg.Infof("Stop %v.", svc.Name)
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
	buildCmd := exec.Command("make", "install", "-f", "./"+svc.Name+"/Makefile")
	buildCmd.Stderr = os.Stderr
	buildCmd.Stdout = os.Stdout
	err := buildCmd.Start()
	if err != nil {
		return fmt.Errorf("can't start build: %v", err)
	}
	err = buildCmd.Wait()
	if err != nil {
		return fmt.Errorf("can't finish build: %v", err)
	}
	return nil
}

func (svc *service) Stop() {
	if svc.IsRunning {
		svc.currentServiceChannel <- KILL_SIGNAL
		svc.syncMutex.Lock()
		svc.syncMutex.Unlock()
	}
}

func (svc *service) SetupService(args ...string) error {
	if svc.currentExternalCmd != nil || svc.currentServiceChannel != nil {
		return fmt.Errorf("service %v already in use", svc.Name)
	} else {
		err := svc.buildService()
		if err != nil {
			return err
		}
		gopath := os.Getenv("GOPATH")
		if gopath == "" {
			return fmt.Errorf("GOPATH is empty")
		} else {
			svc.currentServiceChannel = make(chan interface{})
			svc.Args = append(svc.Args, args...)
			runArgs := []string{}
			runArgs = append(runArgs, svc.Args...)

			cmd := exec.Command(filepath.Join(gopath, "bin", svc.Name), runArgs...)
			svc.currentExternalCmd = cmd
			return nil
		}
	}
}
