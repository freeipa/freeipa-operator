package helpers

import (
	"os"
	"strings"
)

const (
	ENV_TEST_DISABLE_LIST = "TEST_DISABLE_LIST"
)

func IsTestGroupDisabled(groupName string) bool {
	list := os.Getenv(ENV_TEST_DISABLE_LIST)
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
