package command

import (
	"strings"
)

// Todo : add check for existing capabilities provided by the back-end

func getDockerCapabilitiesArray(capabilities string) []string {
	if capabilities == "" {
		return []string{}
	}
	return strings.Split(capabilities, ",")
}
