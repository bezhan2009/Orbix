package utils

import (
	"fmt"
	"goCmd/cmd/commands/commandsWithSignaiture/CF"
)

func DFUtil(commandArgs []string) {
	isSuccess, err := CF.CreateFolder(commandArgs)
	if err != nil {
		fmt.Println("Error deleting folder:", err)
		return
	}

	if isSuccess {
		fmt.Printf("the folder has been deleted: %s\n", commandArgs[0])
	}
}
