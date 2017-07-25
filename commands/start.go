package commands

import (
	"github.com/vetcher/easystarter/backend"
)

type StartCommand struct {
	allFlag bool
}

func (c StartCommand) Validate(args ...string) error {
	if len(args) > 0 {
		c.allFlag = args[0] == ALL
		return nil
	}
	return AtLeastOneArgumentErr
}

func (c StartCommand) Exec(args ...string) error {
	if c.allFlag {
		backend.StartAllServices(args[1:]...)
	} else {
		svc := backend.GetService(args[0], args[1:]...)
		if svc != nil {
			go svc.Start()
		}
	}
	return nil
}
