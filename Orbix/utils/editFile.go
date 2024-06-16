package utils

import (
	"fmt"
	"goCmd/commands/commandsWithSignaiture/Edit"
)

func EditFileUtil(commandArgs []string) {
	if len(commandArgs) < 1 {
		fmt.Println("Usage: edit <file>")
		return
	}
	if err := Edit.File(commandArgs[0]); err != nil {
		fmt.Println(err)
	}
}
