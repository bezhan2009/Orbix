package utils

import (
	"fmt"
	"goCmd/cmd/commands/commandsWithSignaiture/Remove"
	"goCmd/internal/logger"
)

func RemoveFileUtil(commandArgs []string, command string, user, dir string) {
	_, err := Remove.File(commandArgs, command)
	if err != nil {
		logger.Commands(command, false, commandArgs, user, dir)
		fmt.Println(err)
	} else {
		logger.Commands(command, true, commandArgs, user, dir)
	}
}
