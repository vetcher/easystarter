package backend

import (
	"fmt"
	"os"
	"os/exec"
	"sync"
	"time"

	"path/filepath"

	"syscall"

	"github.com/kpango/glg"
)

const (
	OK_SIGNAL   = "ok"
	STOP_SIGNAL = "stop"
	KILL_SIGNAL = "kill"
)

type service struct {
	// Имя сервиса
	Name string
	// Аргументы, которые будут переданы в сервис как параметры командной строки
	Args []string
	// Канал, принимающий сообщения для действий
	serviceSignalChannel chan string
	serviceErrorChannel  chan error
	externalCmd          *exec.Cmd
	// Отметка о старте
	StartTime time.Time
	Target    string
	IsRunning bool
	// Не даёт запустить новый сервис, если старый еще не закончил работу
	// Не дает перезаписывать каналы и `externalCmd` поля
	syncMutex sync.Mutex
}

var allServices map[string]*service

func (svc *service) SetupService(args ...string) error {
	if svc.externalCmd != nil || svc.serviceSignalChannel != nil {
		return fmt.Errorf("service %v already in use", svc.Name)
	} else {
		svc.syncMutex.Lock()
		err := svc.buildService()
		if err != nil {
			return err
		}
		gopath := os.Getenv("GOPATH")
		if gopath == "" {
			return fmt.Errorf("GOPATH is empty")
		}
		// Remember signal channel
		svc.serviceSignalChannel = make(chan string)
		svc.serviceErrorChannel = make(chan error)
		svc.Args = append(svc.Args, args...)

		svc.externalCmd = exec.Command(filepath.Join(gopath, "bin", svc.Name), svc.Args...)
		return nil
	}
}

func (svc *service) buildService() error {
	buildCmd := exec.Command("make", "install", "-f", svc.Target)
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

func (svc *service) Start() {
	err := svc.logInit()
	if err != nil {
		glg.Warn(err)
	}
	err = svc.startService()
	if err != nil {
		glg.Error(err)
	}
	// Now service really started
	svc.handleSignals()
	// Self cleaning because we are not pigs
	svc.cleanService()
	glg.Infof("Stop %v.", svc.Name)
}

func (svc *service) logInit() error {
	out, err := os.Create(fmt.Sprintf("./logs/%v.log", svc.Name))
	// Init log file and all output would write to file
	// If init unsuccessful out will be written to Stdout and Stderr
	if err != nil {
		svc.externalCmd.Stdout = os.Stdout
		svc.externalCmd.Stderr = os.Stderr
		return fmt.Errorf("can't init %v.log file: %v", svc.Name, err)
	} else {
		svc.externalCmd.Stdout = out
		svc.externalCmd.Stderr = out
		return nil
	}
}

func (svc *service) handleSignals() {
	svc.IsRunning = true
	for {
		select {
		case sig := <-svc.serviceSignalChannel:
			switch sig {
			case OK_SIGNAL:
				return
			case KILL_SIGNAL:
				err := svc.externalCmd.Process.Kill()
				if err != nil {
					glg.Errorf("Service %v can't be killed: %v.", svc.Name, err)
				}
				return
			case STOP_SIGNAL:
				err := svc.externalCmd.Process.Signal(syscall.SIGTERM)
				if err != nil {
					glg.Errorf("Service %v can't be stopped: %v.", svc.Name, err)
				}
				return
			}
		case err := <-svc.serviceErrorChannel:
			glg.Errorf("Service %v error:", svc.Name)
			glg.Errorf("%v", err)
			return
		}
	}
}

func (svc *service) waitExecExit() {
	err := svc.externalCmd.Wait()
	if err != nil {
		svc.serviceErrorChannel <- err
	} else {
		svc.serviceSignalChannel <- OK_SIGNAL
	}
}

func (svc *service) startService() error {
	err := svc.externalCmd.Start()
	if err != nil {
		return fmt.Errorf("can't start service %v: %v.", svc.Name, err)
	}
	svc.StartTime = time.Now()
	go svc.waitExecExit()
	glg.Infof("Start %v.", svc.Name, svc.Args)
	return nil
}

func (svc *service) cleanService() {
	close(svc.serviceSignalChannel)
	svc.serviceSignalChannel = nil
	svc.externalCmd = nil
	svc.IsRunning = false
	svc.StartTime = time.Time{}
	svc.syncMutex.Unlock()
	// Now service really stopped
}

func (svc *service) Stop() {
	if svc.IsRunning {
		svc.serviceSignalChannel <- STOP_SIGNAL
	}
}

func (svc *service) Kill() {
	if svc.IsRunning {
		svc.serviceSignalChannel <- KILL_SIGNAL
	}
}

func (svc *service) String() string {
	isRunningStr := "Down"
	if svc.IsRunning {
		isRunningStr = fmt.Sprintf("Up for %v", time.Since(svc.StartTime))
	}
	return fmt.Sprintf("%v\t%v\t%v", svc.Name, isRunningStr, svc.Args)
}
