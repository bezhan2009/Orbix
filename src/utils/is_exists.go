package utils

import (
	"fmt"
	"goCmd/cmd/commands"
	"goCmd/system"
)

func IsExistsUtil(commandArgs []string) bool {
	colors := system.GetColorsMap()

	if len(commandArgs) < 1 {
		fmt.Println(colors["yellow"]("exists <file_name>"))
		return false
	}

	if err := commands.IsExists(commandArgs[0]); err != nil {
		return false
	}

	return true
}
