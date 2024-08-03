package Read

import (
	"fmt"
	utils2 "goCmd/cmd/commands/commandsWithSignaiture/Read/utils"
	"goCmd/internal/logger"
	"goCmd/utils"
)

func File(command string, commandArgs []string, user string, dir string) error {
	if len(commandArgs) < 1 {
		utils.AnimatedPrint(fmt.Sprint("Usage: read <file>"))
		logger.Commands(command, false, commandArgs, user, dir)
		return nil
	}
	nameFileForRead := commandArgs[0]

	dataRead, errReading := utils2.File(nameFileForRead)
	if errReading != nil {
		logger.Commands(command, false, commandArgs, user, dir)
		utils.AnimatedPrint(fmt.Sprint(errReading, "\n"))
		return errReading
	} else {
		logger.Commands(command, true, commandArgs, user, dir)
		fmt.Println(string(dataRead))
		return nil
	}
}
