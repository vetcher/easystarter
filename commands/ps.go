package commands

import (
	"github.com/kpango/glg"
	"github.com/vetcher/easystarter/backend"
)

type PSCommand struct {
}

func (c PSCommand) Exec(args ...string) error {
	if len(args) > 0 {
		glg.Info(backend.ServicesString(args[0]))
	} else {
		glg.Info(backend.ServicesString(""))
	}
	return nil
}
