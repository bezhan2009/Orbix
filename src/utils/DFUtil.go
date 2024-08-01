package utils

import (
	"fmt"
	"goCmd/cmd/commands/commandsWithSignaiture/DF"
)

func DFUtil(commandArgs []string) {
	isSuccess, _ := DF.DeleteFolder(commandArgs)

	if isSuccess {
		fmt.Printf("the folder has been deleted: %s\n", commandArgs[0])
	}
}
