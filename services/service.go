package services

import "time"

const (
	OK_SIGNAL   = "ok"
	STOP_SIGNAL = "stop"
	KILL_SIGNAL = "kill"
)

type Service interface {
	Name() string
	Build() error
	Start() error
	Stop() error
	Kill() error
	Info() *ServiceInfo
	Sync()
	IsRunning() bool
}

type ServiceInfo struct {
	Name        string
	Status      string
	StartupTime time.Time
	Args        []string
}
