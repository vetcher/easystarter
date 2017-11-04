package commands

import (
	"errors"
)

type ExitCommand struct {
}

func (c *ExitCommand) Validate(args ...string) error {
	return nil
}

func (c *ExitCommand) Exec() error {
	return errors.New("exit")
}
