package utils

import (
	"fmt"
	"github.com/fatih/color"
	"goCmd/cmd/commands/commandsWithSignaiture/Edit"
)

func EditFileUtil(commandArgs []string) {
	yellow := color.New(color.FgYellow).SprintFunc()
	if len(commandArgs) < 1 {
		fmt.Println(yellow("Usage: edit <file>"))
		return
	}
	if err := Edit.File(commandArgs[0]); err != nil {
		fmt.Println(err)
	}
}
