package helpers

import (
	"os"
	"strings"
)

const (
	ENV_TESTS_DISABLED = "TESTS_DISABLED"
)

func IsTestGroupDisabled(groupName string) bool {
	list := os.Getenv(ENV_TESTS_DISABLED)
	if list == "" {
		return false
	}
	items := strings.Split(list, ",")
	groupName = strings.TrimSpace(groupName)
	for _, item := range items {
		if strings.TrimSpace(item) == groupName {
			return true
		}
	}
	return false
}
