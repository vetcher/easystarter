package services

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
	String() string
	Sync()
	IsRunning() bool
}
