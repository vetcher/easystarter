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
	"github.com/vetcher/easystarter/backend/util"
)

type goService struct {
	info                 ServiceInfo    // Общая, описательная информация о сервисе
	externalCmd          *exec.Cmd      // Внешняя команда: она является олицетворением сервиса в системе
	serviceSignalChannel chan string    // Канал, пропускающий сообщения для действий
	serviceErrorChannel  chan error     // Канал, пропускающий ошибки исполнения внешней команды
	isRunning            bool           // Метка о состоянии сервиса
	syncMutex            sync.Mutex     // Не дает перезаписывать каналы и `externalCmd` поля, не дает создавать новую cmd, пока старая запущена
	oneInstance          sync.WaitGroup // Для синхронизованной остановки и перезапуска сервисов
	logs                 *os.File       // Файл, куда писать логи
}

func (svc *goService) Name() string {
	return svc.info.Name
}

func (svc *goService) Prepare() error {
	if svc.IsRunning() {
		return fmt.Errorf("%s: already in use", svc.info.Name)
	}
	svc.oneInstance.Wait()
	svc.isRunning = true
	err := svc.prepare()
	if err != nil {
		svc.cleanService()
		return fmt.Errorf("prepare error: %v", err)
	}
	err = svc.logInit()
	if err != nil {
		svc.cleanService()
		return fmt.Errorf("can't init logs: %v", err)
	}
	fmt.Fprintln(svc.logs, svc.Name(), "preparing...")
	return nil
}

func (svc *goService) Build() error {
	fmt.Fprintln(svc.logs, svc.Name(), "building...")
	buildCmd := exec.Command("make", "install", "-f", filepath.Join(svc.info.Target))
	buildCmd.Stderr = svc.logs
	buildCmd.Stdout = svc.logs
	buildCmd.Dir = filepath.Join(svc.info.Dir, svc.info.Name)
	err := buildCmd.Start()
	if err != nil {
		svc.cleanService()
		return fmt.Errorf("can't start build: %v", err)
	}
	err = buildCmd.Wait()
	if err != nil {
		svc.cleanService()
		return fmt.Errorf("can't finish build: %v", err)
	}
	return nil
}

func (svc *goService) Start() error {
	fmt.Fprintln(svc.logs, svc.Name(), "starting...")
	err := svc.startService()
	if err != nil {
		svc.cleanService()
		return fmt.Errorf("%s: can't start service: %v", svc.info.Name, err)
	}
	// Now service really started
	go func() {
		err := svc.handleSignals()
		if err != nil {
			fmt.Fprintf(svc.logs, "%s: handle signals error: %v\n", svc.Name(), err)
		}
		// Self cleaning because we are not pigs
		svc.cleanService()
	}()
	return nil
}

func (svc *goService) logInit() error {
	out, err := os.OpenFile(filepath.Join(util.StartupDir(), "logs", fmt.Sprintf("%s.log", svc.info.Name)), os.O_APPEND|os.O_WRONLY, 0600)
	// Init log file and all output would write to file
	// If init unsuccessful out will be written to Stdout and Stderr
	if err != nil {
		return fmt.Errorf("can't create %s.log file: %v", svc.info.Name, err)
	}
	svc.logs = out
	svc.externalCmd.Stdout = svc.logs
	svc.externalCmd.Stderr = svc.logs
	return nil
}

func (svc *goService) startService() error {
	err := svc.externalCmd.Start()
	if err != nil {
		return fmt.Errorf("can't exec %s: %v", svc.info.Name, err)
	}
	svc.info.StartupTime = time.Now()
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
	for _, arg := range svc.info.Args {
		runArgs = append(runArgs, os.ExpandEnv(arg))
	}

	svc.externalCmd = exec.Command(filepath.Join(gopath, "bin", svc.info.Name), runArgs...)
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
					return fmt.Errorf("can't be killed: %v", err)
				}
			case STOP_SIGNAL:
				err := svc.externalCmd.Process.Signal(syscall.SIGTERM)
				if err != nil {
					return fmt.Errorf("can't be stopped: %v", err)
				}
			}
		case err := <-svc.serviceErrorChannel:
			return fmt.Errorf("service error: %v", err)
		}
	}
}

func (svc *goService) waitExecExit() {
	svc.oneInstance.Add(1)
	defer svc.oneInstance.Done()
	err := svc.externalCmd.Wait()
	if err != nil {
		svc.serviceErrorChannel <- err
	} else {
		svc.serviceSignalChannel <- OK_SIGNAL
	}
}

func (svc *goService) cleanService() {
	fmt.Fprintln(svc.logs, svc.Name(), "cleaning...")
	svc.oneInstance.Done()
	svc.oneInstance.Wait()
	close(svc.serviceSignalChannel)
	close(svc.serviceErrorChannel)
	svc.serviceSignalChannel = nil
	svc.serviceErrorChannel = nil
	svc.logs.Close()
	svc.externalCmd = nil
	svc.info.StartupTime = time.Time{}
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
		Name:        glg.Cyan(svc.info.Name),
		Status:      status,
		Args:        svc.info.Args,
		StartupTime: svc.info.StartupTime,
		Dir:         svc.info.Dir,
		Target:      svc.info.Target,
		Version:     svc.info.Version,
	}
}

func (svc *goService) Sync() error {
	svc.oneInstance.Wait()
	return nil
}

func (svc *goService) IsRunning() bool {
	return svc.isRunning
}
