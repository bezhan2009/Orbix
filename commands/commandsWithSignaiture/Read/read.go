package Read

import (
	"fmt"
	utils2 "goCmd/commands/commandsWithSignaiture/Read/utils"
	"goCmd/debug"
	"goCmd/utils"
)

func File(command string, commandArgs []string) {
	if len(commandArgs) < 1 {
		utils.AnimatedPrint(fmt.Sprint("Использование: read <файл>"))
		return
	}
	nameFileForRead := commandArgs[0]

	dataRead, errReading := utils2.File(nameFileForRead)
	if errReading != nil {
		debug.Commands(command, false)
		utils.AnimatedPrint(fmt.Sprint(errReading, "\n"))
	} else {
		debug.Commands(command, true)
		//_, errWrite := os.Stdout.Write(dataRead)
		utils.AnimatedPrint(string(dataRead))
		//if errWrite != nil {
		//	utils.AnimatedPrint(fmt.Sprint(errWrite, "\n"))
		//}
	}
}
