package Write

import (
	"fmt"
	"goCmd/commands/commandsWithSignaiture/Write/utils"
	"goCmd/debug"
	"strings"
)

func File(command string, commandArgs []string, user string, dir string) error {
	if len(commandArgs) < 2 {
		fmt.Println("Usage: write <file> <data>")
		debug.Commands(command, false, commandArgs, user, dir)
		return nil
	}

	nameFileForWrite := commandArgs[0]
	data := strings.Join(commandArgs[1:], " ")

	if nameFileForWrite == "debug.txt" {
		debug.Commands(command, false, commandArgs, user, dir)
		fmt.Println("PermissionDenied: You cannot write, delete or create a debug.txt file")
		return nil
	}

	errWriting := utils.WriteFile(nameFileForWrite, data+"\n")

	if errWriting != nil {
		debug.Commands(command, false, commandArgs, user, dir)
		fmt.Println(errWriting)
		return errWriting
	} else {
		debug.Commands(command, true, commandArgs, user, dir)
		fmt.Printf("Successfully wrote data to file %s\n", nameFileForWrite)
		return nil
	}
}
