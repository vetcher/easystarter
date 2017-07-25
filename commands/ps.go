package commands

import (
	"github.com/kpango/glg"
	"github.com/vetcher/easystarter/backend"
)

type PSCommand struct {
	allFlag string
}

func (c PSCommand) Validate(args ...string) error {
	if len(args[0]) > 0 {
		c.allFlag = args[0]
	} else {
		c.allFlag = ""
	}
	return nil
}

func (c PSCommand) Exec(args ...string) error {
	glg.Info(backend.ServicesString(c.allFlag))
	return nil
}
