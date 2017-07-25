package commands

import (
	"errors"
)

type exitError struct {
	str string
}

func NewExitError(err string) error {
	return exitError{str:err}
}

func (e exitError) Error() string {
	return e.str
}

type ExitCommand struct {
}

func (c *ExitCommand) Validate(args ...string) error {
	return nil
}

func (c *ExitCommand) Exec(args ...string) error {
	return errors.New("I'm out")
}
