package commands

type Command interface {
	Exec(args ...string) error
}

type EmptyCommand struct {
}

func (c EmptyCommand) Exec(args ...string) error {
	return nil
}
