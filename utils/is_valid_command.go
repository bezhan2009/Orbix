package utils

import (
	"goCmd/system"
)

func ValidCommand(command string, commands []system.Command) bool {
	for _, cmd := range commands {
		if command == cmd.Name {
			return true
		}
	}

	return false
}

func IsValid(command string, commands []string) bool {
	for _, cmd := range commands {
		if command == cmd {
			return true
		}
	}

	return false
}
