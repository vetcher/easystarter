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
	Dir       string
	IsRunning bool
	// Не даёт запустить новый сервис, если старый еще не закончил работу
	// Не дает перезаписывать каналы и `externalCmd` поля
	syncMutex sync.Mutex
	// Для синхронизованной остановки сервисов
	Done sync.WaitGroup
}

var (
	allServices map[string]*service
	STARTUP_DIR string
)

func init() {
	runningDir, err := os.Getwd()
	if err != nil {
		glg.Fatalf("getwd fatal error: %v", err)
	}
	STARTUP_DIR = runningDir
}

func (svc *service) SetupService(args ...string) error {
	if svc.externalCmd != nil || svc.serviceSignalChannel != nil {
		return fmt.Errorf("service %v already in use", svc.Name)
	} else {
		err := svc.buildService()
		if err != nil {
			return err
		}
		gopath := os.Getenv("GOPATH")
		if gopath == "" {
			return fmt.Errorf("GOPATH is empty")
		}
		svc.syncMutex.Lock()
		// Remember signal channel
		svc.serviceSignalChannel = make(chan string)
		svc.serviceErrorChannel = make(chan error)
		svc.Done.Add(1)
		svc.Args = append(svc.Args, args...)
		var runArgs []string
		for _, arg := range svc.Args {
			runArgs = append(runArgs, os.ExpandEnv(arg))
		}

		svc.externalCmd = exec.Command(filepath.Join(gopath, "bin", svc.Name), runArgs...)
		return nil
	}
}

func (svc *service) buildService() error {
	buildCmd := exec.Command("make", "install", "-f", filepath.Join(svc.Target))
	buildCmd.Stderr = os.Stderr
	buildCmd.Stdout = os.Stdout
	buildCmd.Dir = filepath.Join(svc.Dir, svc.Name)
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
	glg.Infof("STOP %v", glg.Yellow(svc.Name))
}

func (svc *service) logInit() error {
	out, err := os.Create(filepath.Join(STARTUP_DIR, "logs", fmt.Sprintf("%s.log", svc.Name)))
	// Init log file and all output would write to file
	// If init unsuccessful out will be written to Stdout and Stderr
	if err != nil {
		svc.externalCmd.Stdout = os.Stdout
		svc.externalCmd.Stderr = os.Stderr
		return fmt.Errorf("can't init %s.log file: %v", svc.Name, err)
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
	glg.Infof("START %s\033[32m with %v\033[39m", glg.Yellow(svc.Name), svc.Args) // \033[32m is Green color, \033[39m is colorless
	return nil
}

func (svc *service) cleanService() {
	close(svc.serviceSignalChannel)
	close(svc.serviceErrorChannel)
	svc.serviceSignalChannel = nil
	svc.serviceErrorChannel = nil
	svc.externalCmd = nil
	svc.StartTime = time.Time{}
	svc.IsRunning = false
	svc.Done.Done()
	svc.syncMutex.Unlock()
	// Now service really stopped
}

func (svc *service) Stop() {
	if svc.IsRunning {
		svc.serviceSignalChannel <- STOP_SIGNAL
	}
}

func (svc *service) SyncStop() {
	if svc.IsRunning {
		svc.serviceSignalChannel <- STOP_SIGNAL
		svc.Done.Wait()
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
	return fmt.Sprintf("%v\t%v\t%v", glg.Yellow(svc.Name), isRunningStr, svc.Args)
}
