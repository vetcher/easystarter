package services

import "time"

const (
	OK_SIGNAL   = "ok"
	STOP_SIGNAL = "stop"
	KILL_SIGNAL = "kill"
)

type Service interface {
	Name() string       // Name of service. By this name manager provide operations.
	Prepare() error     // All what should be done first
	Build() error       // All actions that should happened before run.
	Start() error       // Start service.
	Stop() error        // Normal stop service.
	Kill() error        // Fast stop service.
	Info() *ServiceInfo // General information about current service for `ps` command.
	Sync()              // Sync service with command line. Should be called after Stop/Kill or before Start to prevent desync.
	IsRunning() bool    // Check, is service is idling or working.
}

type ServiceInfo struct {
	Name        string
	Status      string
	StartupTime time.Time
	Args        []string
	Dir         string
	Target      string
}
