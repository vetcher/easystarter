package util

import (
	"errors"
	"strings"
)

func AutoCompleteString(completeThis string, predictFrom []string) (string, int, error) {
	for i, str := range predictFrom {
		if strings.HasPrefix(str, completeThis) {
			return str, i, nil
		}
	}
	return "", -1, errors.New("can't auto-complete")
}

func StringOrEmpty(str string) string {
	if str == "" {
		return "<empty>"
	}
	return str
}

func StrInStrs(str string, strs []string) bool {
	for _, s := range strs {
		if s == str {
			return true
		}
	}
	return false
}

func IndexStrInStrs(str string, strs []string) int {
	for i, s := range strs {
		if s == str {
			return i
		}
	}
	return -1
}