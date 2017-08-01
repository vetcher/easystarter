package util

import (
	"fmt"
	"os"
	"strings"

	"errors"

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

type OverwriteError struct {
	content string
}

func (e *OverwriteError) Error() string {
	return e.content
}

func NewOverwriteError(a string) error {
	return &OverwriteError{
		content: a,
	}
}

func ComposeErrors(errs []error) error {
	if len(errs) > 0 {
		var strs []string
		for _, err := range errs {
			strs = append(strs, err.Error())
		}
		return fmt.Errorf("many errors:\n%v", strings.Join(strs, "\n"))
	}
	return nil
}

func StrInStrs(str string, strs []string) bool {
	for _, s := range strs {
		if s == str {
			return true
		}
	}
	return false
}

func AutoCompleteString(completeThis string, predictFrom []string) (string, int, error) {
	for i, str := range predictFrom {
		if strings.HasPrefix(str, completeThis) {
			return str, i, nil
		}
	}
	return "", -1, errors.New("can't auto-complete")
}
