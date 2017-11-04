package commands

import (
	"errors"

	"github.com/vetcher/easystarter/backend/services"
	"github.com/vetcher/easystarter/client/util"
)

const (
	ALL    = "-all"
	RELOAD = "-reload"
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

func CompleteNames(beforeStrs []string) []string {
	var afterStrs []string
	svcNames := <-services.ServeAllServicesNames()
	for _, name := range beforeStrs {
		completedName, position, err := util.AutoCompleteString(name, svcNames)
		if err != nil {
			afterStrs = append(afterStrs, name)
		} else {
			afterStrs = append(afterStrs, completedName)
			svcNames = append(svcNames[:position], svcNames[position+1:]...) // remove from list
		}
	}
	return afterStrs
}
