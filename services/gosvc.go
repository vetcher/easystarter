package services

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sync"
	"syscall"
	"time"

	"github.com/kpango/glg"
	"github.com/vetcher/easystarter/util"
)

type goService struct {
	// Имя сервиса
	SvcName string
	// Аргументы, которые будут переданы в сервис как параметры командной строки
	Args []string
	// Канал, принимающий сообщения для действий
	serviceSignalChannel chan string
	serviceErrorChannel  chan error
	externalCmd          *exec.Cmd
	// Отметка о старте
	startTime time.Time
	Target    string
	Dir       string
	isRunning bool
	// Не даёт запустить новый сервис, если старый еще не закончил работу
	// Не дает перезаписывать каналы и `externalCmd` поля
	syncMutex sync.Mutex
	// Для синхронизованной остановки и перезапуска сервисов
	oneInstance sync.WaitGroup
}

func (svc *goService) Name() string {
	return svc.SvcName
}

func (svc *goService) Build() error {
	if svc.IsRunning() {
		return fmt.Errorf("service %v already in use", svc.SvcName)
	} else {
		buildCmd := exec.Command("make", "install", "-f", filepath.Join(svc.Target))
		buildCmd.Stderr = os.Stderr
		buildCmd.Stdout = os.Stdout
		buildCmd.Dir = filepath.Join(svc.Dir, svc.SvcName)
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
}

func (svc *goService) Start() error {
	if svc.IsRunning() {
		return fmt.Errorf("service %v already in use", svc.SvcName)
	} else {
		svc.isRunning = true
		err := svc.prepare()
		if err != nil {
			svc.isRunning = false
			return fmt.Errorf("prepare error: %v", err)
		}
		err = svc.logInit()
		if err != nil {
			svc.isRunning = false
			return fmt.Errorf("can't init logs: %v", err)
		}
		err = svc.startService()
		if err != nil {
			svc.isRunning = false
			return fmt.Errorf("can't start service %v: %v", svc.SvcName, err)
		}
		// Now service really started
		go func() {
			svc.handleSignals()
			// Self cleaning because we are not pigs
			svc.cleanService()
		}()
		return nil
	}
}

func (svc *goService) logInit() error {
	out, err := os.Create(filepath.Join(util.StartupDir(), "logs", fmt.Sprintf("%s.log", svc.SvcName)))
	// Init log file and all output would write to file
	// If init unsuccessful out will be written to Stdout and Stderr
	if err != nil {
		svc.externalCmd.Stdout = os.Stdout
		svc.externalCmd.Stderr = os.Stderr
		return fmt.Errorf("can't create %s.log file: %v", svc.SvcName, err)
	} else {
		svc.externalCmd.Stdout = out
		svc.externalCmd.Stderr = out
		return nil
	}
}

func (svc *goService) startService() error {
	err := svc.externalCmd.Start()
	if err != nil {
		return fmt.Errorf("can't start service %v: %v.", svc.SvcName, err)
	}
	svc.startTime = time.Now()
	go svc.waitExecExit()
	return nil
}

func (svc *goService) prepare() error {
	gopath := os.Getenv("GOPATH")
	if gopath == "" {
		return fmt.Errorf("GOPATH is empty")
	}
	svc.syncMutex.Lock()
	svc.serviceSignalChannel = make(chan string)
	svc.serviceErrorChannel = make(chan error)
	svc.oneInstance.Add(1)
	var runArgs []string
	for _, arg := range svc.Args {
		runArgs = append(runArgs, os.ExpandEnv(arg))
	}

	svc.externalCmd = exec.Command(filepath.Join(gopath, "bin", svc.SvcName), runArgs...)
	return nil
}

func (svc *goService) handleSignals() error {
	for {
		select {
		case sig := <-svc.serviceSignalChannel:
			switch sig {
			case OK_SIGNAL:
				return nil
			case KILL_SIGNAL:
				err := svc.externalCmd.Process.Kill()
				if err != nil {
					return fmt.Errorf("service %v can't be killed: %v", svc.SvcName, err)
				}
				return nil
			case STOP_SIGNAL:
				err := svc.externalCmd.Process.Signal(syscall.SIGTERM)
				if err != nil {
					return fmt.Errorf("service %v can't be stopped: %v", svc.SvcName, err)
				}
				return nil
			}
		case err := <-svc.serviceErrorChannel:
			return fmt.Errorf("Service %v error:\n%v", svc.SvcName, err)
		}
	}
}

func (svc *goService) waitExecExit() {
	err := svc.externalCmd.Wait()
	if err != nil {
		//println("err", err.Error(), svc.Name(), svc.IsRunning())
		svc.serviceErrorChannel <- err
	} else {
		//println("ok", svc.Name(), svc.IsRunning())
		svc.serviceSignalChannel <- OK_SIGNAL
	}
}

func (svc *goService) cleanService() {
	close(svc.serviceSignalChannel)
	close(svc.serviceErrorChannel)
	svc.serviceSignalChannel = nil
	svc.serviceErrorChannel = nil
	svc.externalCmd = nil
	svc.startTime = time.Time{}
	svc.oneInstance.Done()
	svc.syncMutex.Unlock()
	svc.isRunning = false
	// Now service really stopped
}

func (svc *goService) Stop() error {
	if svc.IsRunning() {
		svc.serviceSignalChannel <- STOP_SIGNAL
	}
	return nil
}

func (svc *goService) Kill() error {
	if svc.IsRunning() {
		svc.serviceSignalChannel <- KILL_SIGNAL
	}
	return nil
}

func (svc *goService) Info() *ServiceInfo {
	status := "DOWN"
	if svc.IsRunning() {
		status = "UP"
	}
	return &ServiceInfo{
		Name:        glg.Cyan(svc.SvcName),
		Status:      status,
		Args:        svc.Args,
		StartupTime: svc.startTime,
	}
}

func (svc *goService) Sync() {
	svc.oneInstance.Wait()
}

func (svc *goService) IsRunning() bool {
	return svc.isRunning
}
