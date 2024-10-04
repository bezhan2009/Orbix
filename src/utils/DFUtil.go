package utils

import (
	"fmt"
	"github.com/fatih/color"
	"goCmd/cmd/commands"
)

func DFUtil(commandArgs []string) {
	green := color.New(color.FgGreen).SprintFunc()
	isSuccess, _ := commands.DeleteFolder(commandArgs)

	if isSuccess {
		printSuccess := fmt.Sprintf("the folder has been deleted: %s\n", commandArgs[0])
		fmt.Printf(green(printSuccess))
	}
}
