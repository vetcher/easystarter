package commands

import (
	"github.com/kpango/glg"
	"github.com/vetcher/easystarter/backend"
)

type PSCommand struct {
}

func (c PSCommand) Exec(args ...string) {
	glg.Info(backend.ServicesString(args[0]))
}
