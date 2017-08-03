package commands

import (
	"github.com/kpango/glg"
)

type VersionCommand struct {
	Version string
}

func (c *VersionCommand) Validate(args ...string) error {
	return nil
}

func (c *VersionCommand) Exec() error {
	glg.Print(c.Version)
	return nil
}
