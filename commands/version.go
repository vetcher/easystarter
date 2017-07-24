package commands

import (
	"github.com/kpango/glg"
)

type VersionCommand struct {
	Version string
}

func (c VersionCommand) Exec(args ...string) error {
	glg.Print(c.Version)
	return nil
}
