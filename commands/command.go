package commands

type Command interface {
	Exec(args... string)
}
