package Write

import (
	"fmt"
	"goCmd/cmd/commands/commandsWithSignaiture/Write/utils"
	"goCmd/internal/logger"
	"strings"
)

func File(command string, commandArgs []string, user string, dir string) error {
	if len(commandArgs) < 2 {
		fmt.Println("Usage: write <file> <data>")
		logger.Commands(command, false, commandArgs, user, dir)
		return nil
	}

	nameFileForWrite := commandArgs[0]
	data := strings.Join(commandArgs[1:], " ")

	if nameFileForWrite == "debug.txt" {
		logger.Commands(command, false, commandArgs, user, dir)
		fmt.Println("PermissionDenied: You cannot write, delete or create a debug.txt file")
		return nil
	}

	errWriting := utils.WriteFile(nameFileForWrite, data+"\n")

	if errWriting != nil {
		logger.Commands(command, false, commandArgs, user, dir)
		fmt.Println(errWriting)
		return errWriting
	} else {
		logger.Commands(command, true, commandArgs, user, dir)
		fmt.Printf("Successfully wrote data to file %s\n", nameFileForWrite)
		return nil
	}
}
