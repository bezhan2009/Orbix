package Read

import (
	"fmt"
	utils2 "goCmd/cmd/commands/commandsWithSignaiture/Read/utils"
	"goCmd/utils"
)

func File(commandArgs []string) error {
	if len(commandArgs) < 1 {
		utils.AnimatedPrint(fmt.Sprint("Usage: read <file>"))
		return nil
	}
	nameFileForRead := commandArgs[0]

	dataRead, errReading := utils2.File(nameFileForRead)
	if errReading != nil {
		return errReading
	} else {
		fmt.Println(string(dataRead))
		return nil
	}
}
