package utils

import (
	"fmt"
	"goCmd/cmd/commands/commandsWithSignaiture/DF"
)

func DFUtil(commandArgs []string) {
	isSuccess, err := DF.DeleteFolder(commandArgs)
	if err != nil {
		fmt.Println("Error deleting folder:", err)
		return
	}

	if isSuccess {
		fmt.Printf("the folder has been deleted: %s\n", commandArgs[0])
	}
}
