package Read

import (
	"fmt"
	"github.com/fatih/color"
	utils2 "goCmd/cmd/commands/commandsWithSignaiture/Read/utils"
	"goCmd/internal/logger"
	"goCmd/utils"
)

func File(command string, commandArgs []string, user string, dir string) error {
	red := color.New(color.FgRed).SprintFunc()
	if len(commandArgs) < 1 {
		utils.AnimatedPrint(fmt.Sprint("Usage: read <file>"))
		return nil
	}
	nameFileForRead := commandArgs[0]

	dataRead, errReading := utils2.File(nameFileForRead)
	if errReading != nil {
		fmt.Println(red(errReading))
		return errReading
	} else {
		err := logger.Commands(command, true, commandArgs, user, dir)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(string(dataRead))
		return nil
	}
}
