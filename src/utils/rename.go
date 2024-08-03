package utils

import (
	"fmt"
	"goCmd/cmd/commands/commandsWithSignaiture/Rename"
	"goCmd/internal/logger"
)

func RenameFileUtil(commandArgs []string, command string, user, dir string) {
	if err := Rename.Rename(commandArgs, command); err != nil {
		logger.Commands(command, false, commandArgs, user, dir)
		fmt.Println(err)
	} else {
		logger.Commands(command, true, commandArgs, user, dir)
	}
}
