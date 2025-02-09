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

// ValidCommandFast проверяет наличие команды в карте.
func ValidCommandFast(command string, commandsMap map[string]struct{}) bool {
	_, ok := commandsMap[command]
	return ok
}
