package utils

import (
	"fmt"
	"goCmd/cmd/commands/commandsWithSignaiture/CF"
)

func CFUtil(commandArgs []string) {
	isSuccess, err := CF.CreateFolder(commandArgs)
	if err != nil {
		fmt.Println("Error creating folder:", err)
		return
	}
	if isSuccess {
		fmt.Printf("Folder created: %s\n", commandArgs[0])
	}
}
