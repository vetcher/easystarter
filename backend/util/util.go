package util

import (
	"fmt"
	"os"
	"strings"

	"github.com/kpango/glg"
)

var startupDir string

func init() {
	var err error
	startupDir, err = os.Getwd()
	if err != nil {
		glg.Fatalf("getwd fatal error: %v", err)
	}
}

func StartupDir() string {
	return startupDir
}

func ComposeErrors(errs []error) error {
	if len(errs) == 0 {
		return nil
	}
	var strs []string
	for _, err := range errs {
		if err != nil {
			strs = append(strs, err.Error())
		}
	}
	if len(strs) == 1 {
		return fmt.Errorf(strs[0])
	}
	if len(strs) > 0 {
		return fmt.Errorf("many errors:\n%v", strings.Join(strs, "\n"))
	}
	return nil
}
