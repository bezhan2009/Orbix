package utils

import (
	"fmt"
	"github.com/fatih/color"
	"goCmd/cmd/commands/commandsWithSignaiture"
)

func CFUtil(commandArgs []string) {
	red := color.New(color.FgRed).SprintFunc()
	green := color.New(color.FgGreen).SprintFunc()

	isSuccess, err := commandsWithSignaiture.CreateFolder(commandArgs)
	if err != nil {
		fmt.Println(red("Error creating folder:", err))
		return
	}
	if isSuccess {
		fmt.Printf(green(fmt.Sprintf("Folder created: %s\n", commandArgs[0])))
	}
}
