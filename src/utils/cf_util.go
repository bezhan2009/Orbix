package utils

import (
	"fmt"
	"goCmd/cmd/commands"
	"goCmd/system"
)

func CFUtil(commandArgs []string) {
	isSuccess, err := commands.CreateFolder(commandArgs)
	if err != nil {
		fmt.Println(system.Red("Error creating folder:", err))
		return
	}
	if isSuccess {
		fmt.Printf(system.Green(fmt.Sprintf("Folder created: %s\n", commandArgs[0])))
	}
}
