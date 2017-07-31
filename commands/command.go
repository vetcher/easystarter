package commands

import "errors"

const (
	ALL    string = "-all"
	RELOAD string = "-reload"
)

var AtLeastOneArgumentErr = errors.New("at least one argument should be specified")

type Command interface {
	Exec() error
	// always should run before call Exec
	Validate(args ...string) error
}

type EmptyCommand struct {
}

func (c *EmptyCommand) Validate(args ...string) error {
	return nil
}

func (c *EmptyCommand) Exec() error {
	return nil
}
