package debug

import (
	"goCmd/commands/Write"
	"os"
)

func Commands(command string, isSuccess bool) error {
	_, errOpenFile := os.Open("debug.txt")
	var errFile error

	if errOpenFile != nil {
		_, errFile = os.Create("debug.txt")
		if errFile != nil {
			return errFile
		}
	}

	var dataCommands string

	dataCommands += "Command: "
	dataCommands += command
	dataCommands += " ; isSuccess: "
	if isSuccess == true {
		dataCommands += "True;\n"
	} else {
		dataCommands += "False;\n"
	}

	errWriteDebug := Write.File("debug.txt", dataCommands)

	if errWriteDebug != nil {
		return errWriteDebug
	}

	return errWriteDebug
}
