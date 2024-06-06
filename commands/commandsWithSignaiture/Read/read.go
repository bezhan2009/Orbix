package Read

import (
	"fmt"
	utils2 "goCmd/commands/commandsWithSignaiture/Read/utils"
	"goCmd/debug"
	"os"
)

func File(command string, commandArgs []string) {
	if len(commandArgs) < 1 {
		fmt.Println("Использование: read <файл>")
		return
	}
	nameFileForRead := commandArgs[0]

	dataRead, errReading := utils2.File(nameFileForRead)
	if errReading != nil {
		debug.Commands(command, false)
		fmt.Println(errReading)
	} else {
		debug.Commands(command, true)
		_, errWrite := os.Stdout.Write(dataRead)
		if errWrite != nil {
			fmt.Println(errWrite)
		}
	}
}
