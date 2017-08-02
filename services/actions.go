package services

import (
	"fmt"
)

type Action func(Service) error

type Step struct {
	Name string
	Do   Action
}

func CallOneByOne(svc Service, states ...Step) error {
	for _, state := range states {
		err := state.Do(svc)
		if err != nil {
			return fmt.Errorf("%s error: %v", state.Name, err)
		}
	}
	return nil
}
